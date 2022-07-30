# Grape

## A simple API Server used for CRUD interaction with Database & Kubernetes API

This microservice project is a **Part 2** of **'Web-Based Project Management System'**. 
 * [Part 1 'WEB Server'](https://github.com/YushchenkoAndrew/mortis-grimreaper)
 * [Part 2 'API Server'](https://github.com/YushchenkoAndrew/grape)
 * [Part 3 'File Server'](https://github.com/YushchenkoAndrew/void)
 * [Part 4 'bot'](https://github.com/YushchenkoAndrew/botodachi)

![System](/img/System.jpg)

So I guess you wondering, why I called it *'Grape'*. Ehh, you know, 'go' + 'api' => 'gapi' ~> 'grape'. Am I the only one who see this. Please say that Im not !!

Anyway, this project is a middle man between Database & Kubernetes API and the rest of the microservice system **'Web-Based Project Management System'**.


Basic summary:
* **[CRUD]** /api/file - DB logic for **File Table**
* **[CRUD]** /api/link - DB logic for **Link Table** & the rest of the system
* **[CRUD]** /api/project - DB logic for **Project Table**
* **[CRUD]** /api/subscription - DB logic for **Subscription Table** *(this handler requires cron logic from **bot**)*
* **[CRUD]** /api/pattern - DB logic for **Pattern Table** *[In progress]*
* **[CRUD]** /api/world - DB logic for **World Table** *[In progress]*
* **[CRUD]** /api/info/{sum/range} - DB logic for **Info Table** *[In progress]*
* **[POST]** /api/bot/redis - Send directly command to Redis, accessed by **bot**
* **[CRUD]** /api/k3s/deployment - A middleman between system and Kubernetes API **Deployments**
* **[CRUD]** /api/k3s/service - A middleman between system and Kubernetes API **Service**
* **[CRUD]** /api/k3s/ingress - A middleman between system and Kubernetes API **Ingress**
* **[CRUD]** /api/k3s/namespace - A middleman between system and Kubernetes API **Namespace**
* **[CRUD]** /api/k3s/pods - A middleman between system and Kubernetes API **Pods** *[In progress]*
* **[CRUD]** /api/k3s/pods/metrics - A middleman between system and Kubernetes API **Pod Metrics** and DB logic for **Metrics Table**


*(More information about routes you can find with swagger docs: http://127.0.0.1:31337/api/swagger/index.html)*

### Models
* **file** - Contains basic info about file which is used by **File Server**
* **project** - Contains basic info about project which is used by **WEB Server**
* **link** - Contains basic info about link which is used by **WEB Server**
* **geo_ip_blocks** & **geo_ip_locations**  - GeoIP2 tables model
* **subscription** - Contains basic info about subscription which is used by **bot**
* **metrics** - Contains basic info about pods metrics which is used by **WEB Server**


## Diagram
![Diagram](/img/API.jpg)

## DB table relations
![DB table relations](/img/Relations.jpeg)

## How to use this project
1. Clone this repo
2. Copy .env.template to .env
```
cp .env.template .env
```
3. Configure a .env
4. Download GeoLite2 csv files
5. Create ./tmp dir
```
mkdir ./tmp
```
6. Copy GeoLite2-Country-Blocks.csv & GeoLite2-Country-Locations-en.csv into ./tmp dir
4. Start the project
```
make dev
```


Now you can check if api is started by sending a **GET** request to http://localhost:31337/api/ping

Swagger docs about routes you can find at http://127.0.0.1:31337/api/swagger/index.html

## Found a bug ?
Found something strange or just want to improve this project, just send a PR.

## What's not implemented yet / Known issues
- [ ] Project test coverage 
  - [x] file
  - [x] link
  - [x] metrics
  - [x] project
  - [ ] subscription
  - [ ] other models....
- [ ] Finish impl of *[In progress]* logic