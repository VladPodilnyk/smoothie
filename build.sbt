val V = new {
  val scalatest     = "3.1.1"
  val scalacheck    = "1.14.3"
  val zio           = "1.0.0-RC18-2"
  val zioCats       = "2.0.0.0-RC13"
  val kindProjector = "0.11.0"
  val jedis         = "3.2.0"
  val cats          = "2.2.0-M1"
}

val D = new {
  val zio           = "dev.zio" %% "zio" % V.zio
  val zioCats       = "dev.zio" %% "zio-interop-cats" % V.zioCats
  val catsCore      = "org.typelevel" %% "cats-core" % V.cats
  val scalatest     = "org.scalatest" %% "scalatest" % V.scalatest
  val scalacheck    = "org.scalacheck" %% "scalacheck" % V.scalacheck
  val kindProjector = "org.typelevel" % "kind-projector" % V.kindProjector cross CrossVersion.full
  val jedis         = "redis.clients" % "jedis" % V.jedis
}

inThisBuild(
  Seq(
    scalaVersion := "2.13.2",
    version := "0.0.1-SNAPSHOT"
  )
)

lazy val root = project
  .in(file("."))
  .settings(
    name := "smoothie",
    libraryDependencies ++= Seq(
      D.zio,
      D.zioCats,
      D.catsCore,
      D.jedis,
      D.scalacheck % Test,
      D.scalatest % Test
    ),
    addCompilerPlugin(D.kindProjector)
  )
