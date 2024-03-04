# Changelog

## [1.1.0](https://github.com/soockee/ssr-go/compare/v1.0.16...v1.1.0) (2024-03-04)


### Features

* enabled HTTPS and update log messages ([795e048](https://github.com/soockee/ssr-go/commit/795e0487071e3e18a306e5439d5e0384ec4eac4f))

## [1.0.16](https://github.com/soockee/ssr-go/compare/v1.0.15...v1.0.16) (2024-03-04)


### Bug Fixes

* fix dockerfile apt install ([93cf8c0](https://github.com/soockee/ssr-go/commit/93cf8c00953dc19739bdc1f03930d122ff2cdca5))

## [1.0.15](https://github.com/soockee/ssr-go/compare/v1.0.14...v1.0.15) (2024-03-04)


### Bug Fixes

* fix dockerfile apt install ([1364b4b](https://github.com/soockee/ssr-go/commit/1364b4b71f37af06f1d90d8bae21c69eeaeca397))

## [1.0.14](https://github.com/soockee/ssr-go/compare/v1.0.13...v1.0.14) (2024-03-04)


### Bug Fixes

* switch to slog for logging ([85adc1d](https://github.com/soockee/ssr-go/commit/85adc1d0c9df1999f4dda1ee7ccb22e84bcbffec))

## [1.0.13](https://github.com/soockee/ssr-go/compare/v1.0.12...v1.0.13) (2024-03-04)


### Bug Fixes

* add flag ([35a8ab5](https://github.com/soockee/ssr-go/commit/35a8ab59e106aa76936d91dfbf3ad3baa2ce4158))
* refactor autocert ([b29e532](https://github.com/soockee/ssr-go/commit/b29e532f5fed207b16571c4441f61bb9af23ba05))
* remove certmanger pointer ([15e8b55](https://github.com/soockee/ssr-go/commit/15e8b555bbe347463c9be067fad8d762d45b0758))
* remove listenAddr default to https/http ([c7e0fe6](https://github.com/soockee/ssr-go/commit/c7e0fe6694fe69abfac80c0815ed95ce775e2788))

## [1.0.12](https://github.com/soockee/ssr-go/compare/v1.0.11...v1.0.12) (2024-03-04)


### Bug Fixes

* add cache dir as function ([6e6f9f5](https://github.com/soockee/ssr-go/commit/6e6f9f595933fbb3282d9caef7a129e6311fc5e3))
* remove secrets from build and prepare deploy stem ([9b29dbb](https://github.com/soockee/ssr-go/commit/9b29dbb44221eb0f1d32b570329d335d6199ce01))

## [1.0.11](https://github.com/soockee/ssr-go/compare/v1.0.10...v1.0.11) (2024-03-04)


### Bug Fixes

* change from secrets to build-args to get them into container ([eae05dd](https://github.com/soockee/ssr-go/commit/eae05ddaf43060d34c31ab31aaf1cd2628a62538))
* remove redirectHandler, because certMang redirects by default ([02cb00c](https://github.com/soockee/ssr-go/commit/02cb00c07d4af5089fdb46c91d2ceec97ba332e1))

## [1.0.10](https://github.com/soockee/ssr-go/compare/v1.0.9...v1.0.10) (2024-03-04)


### Bug Fixes

* change base image for debugging ([8add288](https://github.com/soockee/ssr-go/commit/8add28841fbdb35a62cbba0b443e71a7c57cc2f2))
* refactor https ([5d2d147](https://github.com/soockee/ssr-go/commit/5d2d147227cb75aa27ce9908da1e84e8c60f60d3))

## [1.0.9](https://github.com/soockee/ssr-go/compare/v1.0.8...v1.0.9) (2024-03-03)


### Bug Fixes

* add https ([23e4e36](https://github.com/soockee/ssr-go/commit/23e4e3644cb2f0d08a690c570064f124638ad0c9))

## [1.0.8](https://github.com/soockee/ssr-go/compare/v1.0.7...v1.0.8) (2024-03-03)


### Bug Fixes

* use templ build before go build ([01d2e67](https://github.com/soockee/ssr-go/commit/01d2e674f7e6c35bb988e0825eda11243debb2fb))

## [1.0.7](https://github.com/soockee/ssr-go/compare/v1.0.6...v1.0.7) (2024-03-03)


### Bug Fixes

* add buildx setup step ([d52148d](https://github.com/soockee/ssr-go/commit/d52148d1418a4c1184e9a510cfe5bea2411358dd))

## [1.0.6](https://github.com/soockee/ssr-go/compare/v1.0.5...v1.0.6) (2024-03-03)


### Bug Fixes

* remove go build args ([d015a7c](https://github.com/soockee/ssr-go/commit/d015a7c33c2cf528b190db7a2bdff0b86f9ac03f))
* remove nested module ([0af9f1c](https://github.com/soockee/ssr-go/commit/0af9f1cdcdd2268ddf716072db5b1e22b8365008))

## [1.0.5](https://github.com/soockee/ssr-go/compare/v1.0.4...v1.0.5) (2024-03-03)


### Bug Fixes

* add components ([65a5c56](https://github.com/soockee/ssr-go/commit/65a5c56f5dc4c38ad53b1a310ff69ab24bbacc72))

## [1.0.4](https://github.com/soockee/ssr-go/compare/v1.0.3...v1.0.4) (2024-03-03)


### Bug Fixes

* add fun wip animation ([8d105bd](https://github.com/soockee/ssr-go/commit/8d105bd33ef720107aab81e5c27ce020883624cc))
* add multiplatform built ([7332bb9](https://github.com/soockee/ssr-go/commit/7332bb9230b185dda7b9b6c6bf8e610559f0c84d))
* reduce required go version ([0fe5af0](https://github.com/soockee/ssr-go/commit/0fe5af0454c87933d1647988ec0294f3cc66d83b))

## [1.0.3](https://github.com/soockee/ssr-go/compare/v1.0.2...v1.0.3) (2024-01-15)


### Bug Fixes

* change exposed port ([49c43ff](https://github.com/soockee/ssr-go/commit/49c43ff0f6ae1b8d39f72120e28a853d27c97684))

## [1.0.2](https://github.com/soockee/ssr-go/compare/v1.0.1...v1.0.2) (2024-01-15)


### Bug Fixes

* add templ in build chain ([d685ee0](https://github.com/soockee/ssr-go/commit/d685ee016c794329dd2cd4edd35d0fa7fe0cbd00))

## [1.0.1](https://github.com/soockee/ssr-go/compare/v1.0.0...v1.0.1) (2024-01-14)


### Bug Fixes

* add image publish workflow ([8fb133b](https://github.com/soockee/ssr-go/commit/8fb133b9107508f246b3056185974737c83cdb68))

## 1.0.0 (2024-01-14)


### Features

* add homepage ([8a001da](https://github.com/soockee/ssr-go/commit/8a001dac9f51f19b2417a19819d9252e7ae7193e))


### Bug Fixes

* change readme encoding ([62c2314](https://github.com/soockee/ssr-go/commit/62c2314686c152c62f05f0d7ede4be9fdcf382fa))
* dockerfile ([f036b43](https://github.com/soockee/ssr-go/commit/f036b4316b89060094523f633e1522d750b9c184))
* export component in templates ([d4b7754](https://github.com/soockee/ssr-go/commit/d4b77542a5a09541c519419022630004dba89726))
* initial commit ([42cc89d](https://github.com/soockee/ssr-go/commit/42cc89d9b6c4e5adb78ce649c2f286bb31723d6a))
* remove arch specification for go build ([39d49d0](https://github.com/soockee/ssr-go/commit/39d49d0ed45799435e87af7205a4e59933ade98f))
* replace commit linter ([61bc859](https://github.com/soockee/ssr-go/commit/61bc8593816cee4b67cadb2e6dbed1b4d8cc0a4d))
* templ integration ([b8c1640](https://github.com/soockee/ssr-go/commit/b8c1640840f9b84adc45281e59e8975a54672258))
