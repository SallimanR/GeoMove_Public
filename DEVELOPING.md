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
for building web version:
```sh
pnpm run build
```
for Android version
```sh
pnpm run setup:android
```
```sh
pnpm run build:android
```
for IOS version
```sh
pnpm run setup:ios
```
```sh
pnpm run build:ios
```
