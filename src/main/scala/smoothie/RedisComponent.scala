package smoothie

import redis.clients.jedis.{Jedis, JedisPool}
import zio.{Managed, Task, UIO}

trait RedisComponent {
  def eval[A](f: Jedis => Task[A]): Task[A]
}

object RedisComponent {
  def mkResource(cfg: RedisConfig): Managed[Throwable, Jedis] = Managed.make(Task(new JedisPool(cfg.url).getResource))(a => UIO(a.close()))
}
