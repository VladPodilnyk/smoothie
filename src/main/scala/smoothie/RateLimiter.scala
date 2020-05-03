package smoothie

import zio.Task

trait RateLimiter {
  def start[A](key: String)(f: Task[A]): Task[Unit]
}

object RateLimiter {
  // acuire all necessary resources and create rate limiter
  def mkRateLimiter = ???

  private final class RateLimiterImpl(
    redisComponent: RedisComponent
  ) extends RateLimiter {
    override def start[A](key: String)(f: Task[A]): Task[Unit] = ???


    private[this] def inc(key: String) = {
      redisComponent.eval {
        jedisClient =>
          Task(jedisClient.pipelined()).bracketAuto(???)
      }
    }
    private[this] def dec(key: String) = ???
  }
}
