# GeoMove

Открытая гео-платформа для поиска водителей эвакуаторов и грузоперевозок.

[geomove.online](https://geomove.online) | [driver.geomove.online](https://driver.geomove.online)

Приложения для Android: **[перейти на страницу для скачивания](https://github.com/SallimanR/GeoMove_Public/releases/tag/pre-release_1)** -> скачать apk приложения для пользователей/водителей -> найти apk в загрузках браузера или в папке Downloads -> запустить apk

IOS версии нельзя скачать через GitHub из-за ограничений установленных Apple и IOS

Авторизация через Telegram может не работать без VPN из-за блокировок

## Создание и отслеживание заказа

<p align="center">
  <img src="./docs/screenshots/geomove_customer_order_create.png" alt="Пользователь создаёт заказ" width="400">
  <br>
  Пользователь выбирает локацию отправки и прибытия, создаёт заказ
</p>

<p align="center">
  <img src="./docs/screenshots/geomove_customer_order_driver_app.png" alt="Заказ виден в приложении водителя" width="400">
  <br>
  Водитель видит поступивший заказ в своём приложении
</p>

<p align="center">
  <img src="./docs/screenshots/geomove_customer_order_accepted_driver_app.png" alt="Водитель принимает заказ" width="400">
  <br>
  Водитель принимает заказ — статус меняется в реальном времени
</p>

<p align="center">
  <img src="./docs/screenshots/geomove_customer_order_accepted.png" alt="Пользователь видит, что заказ принят" width="400">
  <br>
  Пользователь получает уведомление и видит обновлённый статус заказа
</p>


<p align="center">
  <img src="./docs/screenshots/freely_available_driver_edit.png" alt="Пользователь создаёт заказ" width="400">
  <br>
  Водитель создаёт заявку на свободный эвакуатор
</p>

<p align="center">
  <img src="./docs/screenshots/freely_available_driver.png" alt="Пользователь создаёт заказ" width="400">
  <br>
  Созданная заявка на свободный эвакуатор
</p>


---

### Библиотеки для использования в сторонних проектах
- #### [Maps npm package](https://www.npmjs.com/package/@geomove/maps) | [docs](./frontend/packages/maps/README.md)
- #### [Geo utilities npm package](https://www.npmjs.com/package/@geomove/go) | [docs](./frontend/packages/geo/README.md)

### Geo API:
- https://geomove.online/style/style/style.json — стиль карт (style.json)
- https://geomove.online/tiles — PMTiles API (тайлы карт)
- https://geomove.online/geocoding — геопоиск и обратная геокодировка
- https://geomove.online/routing — построение маршрутов и map matching

## Лицензия

Проект использует двойное лицензирование:

| Компонент | Лицензия |
|---|---|
| Backend, бизнес-логика, домены | [![AGPL-3.0](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](./LICENSE) |
| [Карты (@geomove/maps)](./frontend/packages/maps) | [![MIT](https://img.shields.io/badge/License-MIT-green.svg)](./frontend/packages/maps/LICENSE) |
| [Гео-утилиты (@geomove/geo)](./frontend/packages/geo) | [![MIT](https://img.shields.io/badge/License-MIT-green.svg)](./frontend/packages/geo/LICENSE) |

Карты и гео-утилиты распространяются под MIT — используйте в любых проектах, включая коммерческие.

---

### Для разработки и участия:
- [DEVELOPING.md](./DEVELOPING.md)
- [Документация](./docs/)
