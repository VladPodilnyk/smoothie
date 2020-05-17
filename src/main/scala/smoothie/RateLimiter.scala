package smoothie

import java.time.ZonedDateTime

import cats.effect.Resource
import zio.Task
import zio.interop.catz._

import scala.jdk.CollectionConverters._

trait RateLimiter {
  def execWithLimit[A](key: String)(f: Task[A]): Task[Unit]
}

object RateLimiter {
  // acuire all necessary resources and create rate limiter
  // TODO: config should be passed via ZLayer
  def mkRateLimiterR(rate: Rate, config: RedisConfig): Resource[Task, RateLimiter] = {
    Resource.make(RedisComponent.mkResource(config).map(new RateLimiterImpl(_, rate)))(_ => Task.unit)
  }

  /** For testing purposes only */
  private[smoothie] def mkRateLimiterF(rate: Rate, config: RedisConfig): Task[RateLimiter] = {
    RedisComponent
      .mkResource(config)
      .map(new RateLimiterImpl(_, rate))
  }

  private final class RateLimiterImpl(
    redisComponent: RedisComponent,
    rate: Rate
  ) extends RateLimiter {
    override def execWithLimit[A](key: String)(f: Task[A]): Task[Unit] = {
      for {
        counter <- get(key)
        _       <- inc(key)
        _       <- Task.when(counter + 1 > rate.n)(Task.fail(new Throwable("You request was rate limited.")))
        _       <- f
      } yield ()
    }

    private[this] def inc(key: String): Task[Unit] = {
      redisComponent.eval {
        jedisClient =>
          Task(jedisClient.pipelined()).bracketAuto {
            pipe =>
              Task {
                val luaScript =
                  """
                    |if tonumber(redis.call('incr', KEYS[1])) == 1 then
                    | return redis.call('expireAt', KEYS[1], ARGV[1])
                    |else
                    | return 0
                    |end
                    |""".stripMargin

                pipe.eval(
                  luaScript,
                  List(key).asJava,
                  List(ZonedDateTime.now().toInstant.plusSeconds(rate.t.toSeconds).getEpochSecond.toString).asJava
                )
              }.unit
          }
      }
    }

    private[this] def get(key: String): Task[Int] = {
      redisComponent.eval {
        jedisClient =>
          Task(jedisClient.get(key))
      }.map(s => Option(s).map(_.toInt).getOrElse(0))
    }
  }
}
