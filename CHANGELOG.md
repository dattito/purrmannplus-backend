# Changelog

<!--
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
-->

## [Unreleased]

## [v0.5.3] - 2021-10-21

### Added

- Ability to specify CORS origins

## [v0.5.3] - 2021-10-20

### Added

- More comments in the code
- Possibility to just login for this session

### Changed

- Changed default version to ""
- Changed response of login with cookie
- Sending 401 if jwt is missing
- Now exposes port in Dockerfile

## [v0.5.2] - 2021-10-13

### Added

- Route to check if the user is logged in

## [v0.5.1] - 2021-10-11

### Added

- Account can now logout
- About Route
- Version for the about route gets passed by the github action

### Changed

- Changed `auth_id` and `auth_pw` to `username` and `password`

### Fixed

- While creating a new account, a user error was referenced as an internal error

## [v0.5.0] - 2021-10-06

### Added

- More debug logging messages
- Ability to disable the api (for example, if it should only use the scheduler
  for updating)

### Fixed

- signal messages can't sometimes send because of utf-8 encoding error

## [v0.4.1] - 2021-10-06

### Changed

- Disabled logging the path
- Changed the logging directory

## [v0.4.0] - 2021-10-05

### Added

- A cookie can be linked to a domain
- Better Logging with different Levels

### Fixed

- Wrong version-number in CHANGELOG.md (v0.2 instead of v0.2.0)

## [v0.3.0] - 2021-10-03

### Changed

- Renamed imports from `/datti-to/` to `/dattito/`
- Phone-Number confirmation message can only send if it wasn't added yet

### Fixed

- Wrong gorm annotations in structs
- Wrong version-number in CHANGELOG.md (v0.1 instead of v0.2)

## [v0.2.0] - 2021-10-02

### Added

- JWT Authentication token can be stored in cookie and read from cookie

## [v0.1.3] - 2021-09-23

### Fixed

- Scheduler was not activated

## [v0.1.2] - 2021-09-22

### Changed

- Docker WORKDIR changed to /data

## [v0.1.1] - 2021-09-22

### Fixed

- Wrong expression in github docker push workflow

## [v0.1] - 2021-09-22

### Added

- Add and delete accounts
- Add phonenumbers to accounts
- (Un)register to the substitution updater which automatically sends updates on
  new substitutions

[Unreleased]: https://github.com/Dattito/purrmannplus-backend/compare/v0.5.4...dev
[v0.5.4]: https://github.com/Dattito/purrmannplus-backend/compare/v0.5.3...v0.5.4
[v0.5.3]: https://github.com/Dattito/purrmannplus-backend/compare/v0.5.2...v0.5.3
[v0.5.2]: https://github.com/Dattito/purrmannplus-backend/compare/v0.5.1...v0.5.2
[v0.5.1]: https://github.com/Dattito/purrmannplus-backend/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/Dattito/purrmannplus-backend/compare/v0.4.1...v0.5.0
[v0.4.1]: https://github.com/Dattito/purrmannplus-backend/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/Dattito/purrmannplus-backend/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/Dattito/purrmannplus-backend/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/Dattito/purrmannplus-backend/compare/v0.1.3...v0.2.0
[v0.1.3]: https://github.com/Dattito/purrmannplus-backend/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/Dattito/purrmannplus-backend/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/Dattito/purrmannplus-backend/compare/v0.1...v0.1.1
[v0.1]: https://github.com/Dattito/purrmannplus-backend/releases/tag/v0.1
