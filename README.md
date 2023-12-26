## Overview

This repository is a template for a Todo application in the Golang that can run on any platform as long as Docker is running.

MySQL is used as the data store.

## Endpoints

Method | Path                                                        | Description             
---  |-------------------------------------------------------------|-------------------------
GET  | `/`                                                         | Main page               |
GET  | `/banners`                                                  | Get all banners         |
GET  | `/banners/{id}`                                             | Get banner by id        |
POST | `/add_banner_to_slot/?banner={bannerId}&slot={slotId}`      | Add banner to slot      |
POST | `/remove_banner_from_slot/?banner={bannerId}&slot={slotId}` | Remove banner from slot |
GET  | `/get_banner_for_show/?group={groupId}&slot={slotId}`       | Get relevant banner     |

## Usage

```shell
$ cp .env.example .env

$ docker compose up -d --build
```

## Tips

### Connect db container

```shell
$ docker compose exec -it db /bin/bash -c "mysql -uroot -p<PASSWORD>"
```
