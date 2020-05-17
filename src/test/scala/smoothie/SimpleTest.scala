package smoothie

import zio.{Ref, Schedule, Task}
import zio.test.Assertion._
import zio.test._
import zio.interop.catz._

import scala.concurrent.duration._

object SimpleTest extends DefaultRunnableSpec {
  override def spec = suite("SimpleTesSuite")(
    testM("super simple test") {
      val config = RedisConfig("http://localhost:6379")
      val rate   = Rate(5, 5.minute)
      val t = for {
        limiter <- RateLimiter.mkRateLimiterF(rate, config)
        ref     <- Ref.make(0)
        task = for {
          _ <- ref.update(_ + 1)
          c <- ref.get
        } yield c
        _   <- limiter.execWithLimit("some")(task).repeat(Schedule.forever).catchAll(_ => Task.unit)
        res <- ref.get
      } yield res

      assertM(t)(equalTo(5))
    }
  )
}
