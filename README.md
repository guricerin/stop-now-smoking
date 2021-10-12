# stop-now-smoking

喫煙本数を記録し、それを監視し合うSNS。

## run local

```bash
$ docker-compose up
```

## deploy

using heroku

- deploy to heroku

```bash
$ git push heroku main
# or
$ git push heroku <herokuに反映させたいoriginのブランチ名>:main
```

- setup db

```bash
$ heroku pg:psql -a stop-now-smoking
>>> \i ./db/initdb.d/setup.sql
>>> \q
```
