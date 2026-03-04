monorepos to have the right versions together

libmxclient - golang matrix client library
pygomx-module - python binding package
smal - python matrix lib


just run the demobot (docker):

docker compose build demobot
docker compose run --rm demobot smalsetup <matrixid>
docker compose up -d demobot

the bot follows each invite (autojoin) and have two commands:
  !stop         - graceful shutdown
  !echo [text]  - in a DM the bot responds with 'text',
                  in regular rooms it is a reply.

install from source (venv):

  build configuration is done via env vars

  # one of static, shared
  PYGOMX_BUILD_MODE=static

  # one of none, colm, goolm, vodozemac
  PYGOMX_OLM_FLAVOR=colm
  # for colm you need libolm-dev installed
  # vodozemac is not supported yet

  you need go >=1.22 installed

  (create and activate a venv)

  cd pygomx-module
  pip install .

  cd ../smal
  pip install [-e] .

usage:

  cd into an empty dir (you might create one)

  smalsetup <matrixid>
    this command creates a credentials file (.mxpass) in the current dir.
    autopickup by all tools & bots that requires credentials, no further configuration required

  commands:
    mxdiscover --help
    mxwhoami
    mxtoken
    mxaccountinfo
    mxclearaccount
    mxserverinfo
    demobot

matrix room:
  #pygomx:matrix.org
  https://matrix.to/#/#pygomx:matrix.org
