monorepos to have the right versions together

libmxclient - golang matrix client library
pygomx - python binding package
smal - python matrix lib


just run the demobot (docker):

docker compose build demobot
docker compose run --rm demobot smalsetup <matrixid>
docker compose up -d demobot

the bot follows each invite (autojoin) and have two commands:
  !stop         - graceful shutdown
  !echo [text]  - in a DM the bot responds with 'text',
                  in regular rooms it is a reply.


binary/package install:

  pip install pygomx
  pip install mxsmal


install from source / develop (venv):

  build requirements:
    golang 1.25+
      git 2.24+
    gcc 4.6+
      windows: gcc built with binutils 2.37+
    python 3.10+

  (obtain this repository and cd into it)

  build configuration is done via env vars

  # one of static, shared
  PYGOMX_BUILD_MODE=static

  # one of none, colm, goolm, vodozemac
  PYGOMX_OLM_FLAVOR=goolm
  # for colm you need libolm-dev installed
  # vodozemac is not supported yet

  (create and activate a venv)

  cd pygomx
  pip install .

  cd ../mxsmal
  pip install [-e] .

usage:

  cd into an empty dir (you might create one)

  smalsetup <matrixid>
    this command creates a credentials file (.mxpass) in the current dir.
    autopickup by all tools & bots that requires credentials, no further configuration required

  commands:
    mxdiscover --help
    mxpassitem --help
    mxwhoami --help
    mxlogout --help
    mxtoken
    mxaccountinfo
    mxclearaccount
    mxserverinfo
    demobot

matrix room:
  #pygomx:matrix.org
  https://matrix.to/#/#pygomx:matrix.org
