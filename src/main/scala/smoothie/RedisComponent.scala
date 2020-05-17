package smoothie

import redis.clients.jedis.{Jedis, JedisPool}
import zio.{Task, UIO}

trait RedisComponent {
  def eval[A](f: Jedis => Task[A]): Task[A]
}

object RedisComponent {
  def mkResource(cfg: RedisConfig): Task[RedisComponent] = {
    Task(new JedisPool(cfg.url).getResource)
      .bracket(
        j => UIO(j.close()),
        j => Task { new RedisComponent { override def eval[A](f: Jedis => Task[A]): Task[A] = f(j) }}
      )
  }
}
