### Commits

#### Use hooks for **message tags** in commits: \<commit text\> + **category**, **subcategories**
examples:
  - "Add \<feature\> \[backend\]"
  - "Change \<something\> \[backend\]\[db\]"

to setup git-hooks run:
```sh
./scripts/setup-hooks.sh
```


### Continous Integration (CI) setup

add `Gitea` to remote origins

after pushing changes the CI will start, you can see the result in Actions tab in `Gitea`


### Running tests locally
```sh
cd backend/monolith
go test ./...
```


### Building the apps
in each apps' directory
##### for building web version:
```sh
pnpm run build
```
##### for Android version:
```sh
pnpm run setup:android
```
```sh
pnpm run build:android
```

1. open project in Android Studio with `<absolute-path-to-repository>/frontend/apps/<app-name>`
2. sync project with "Gradle build" button
3. run the app

##### for IOS version:
```sh
pnpm run setup:ios
```
```sh
pnpm run build:ios
```

1. Open project in XCode
2. build and run

##### To debug mobile version:
1. launch app on a phone connected to same network as development machine
2. open chromium browser on dev machine
3. in chromium go to `<browser-name>://inspect/#devices`
4. select device where app runs
5. click "inspect" -> window with app screen translation (interactive) will appear
6. open browser devtool

generating `open-api`:
```sh
./scripts/generate_openapi.sh
```

generating `protobuf`:
```sh
./scripts/generate_protobuf.sh
```

generating SQL go code with `sqlc` with compile time checks (guarantees validity of queires and data types):
```sh
cd backend/monolith
sqlc generate
```

SQL migrations:
```sh
cd backend/monoltih
```
for applying migrations:
```sh
dbmate up 
```

for reverting migrations:
```sh
dbmate down
```

for dropping the DB with all data and migrations
```sh
dbmate drop
```
