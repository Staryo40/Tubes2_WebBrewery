# Tugas Besar 2 Strategi Algoritma IF2211 2025

## Little Alchemy 2 Element Recipe Finder with BFS, DFS, and Bidirectional Algorithms
Breadth first search (BFS) is a pathfinding algorithm that focuses on exploring each alternative little by little concurrently. On the other hand, Depth First Search (DFS), also a pathfinding algorithm, focuses on choosing one alternative and going through it until it reaches the end before exploring another alternative. Performance-wise, BFS guarantees in finding the most efficient path towards the goal but is cumbersome on memory-usage, while DFS is less memory hungry and faster, but is not guaranteed to find the most efficient path.  

Both algorithms are implemented in creating this simple web application project to find an element in the game called Little Alchemy 2, which is a game where we combine elements to find other more higher tier element. The goal of the game is to find all elements that are possible in the game. There are 720 elements in total in this game without any extra add-ons and finding the recipes to complete the game can become troublesome. That is why this app helps in finding paths to get to a certain element target. Specifically for this project the algorithm used are BFS, DFS, and Bidirectional search. Users can switch between modes to see the performance difference of each algorithm. In addition, there is a feature to show multiple paths towards a single target element that is implemented with multithreading to speed up the process of searcing.

<h3>Deployed Web App: <a href="https://web-brewery-957948630882.asia-southeast2.run.app/">Little Alchemy 2 Recipe Finder by WebBrewery 2025</a><h3>

<p align="center">
<img src="https://github.com/user-attachments/assets/3cc9188f-c9fc-4175-8420-b07e83401b94" alt="Element Recipe Search with BFS/DFS Algorithm Recording" width="700"/>
</p>
<p align="center">Element Recipe Searching</p>
<p align="center">
<img src="https://github.com/user-attachments/assets/a5ae7e72-5f25-4705-a325-e02547b73288" alt="Little Alchemy 2 Elements Encyclopedia" width="700"/>
</p>
<p align="center">Element Encyclopedia</p>

## Program Requirements
### Backend (Go)
  - Go (Latest 1.x version)

### Frontend (Next.js)
- Node.js
- NPM

### Containerized Setup
- Docker & Docker Compose  

## Compiling the Program
There are two ways to compile the program: run source file OR run docker
### Compiling for running locally from source files
1. Compiling Go Backend API:
```bash
cd src/backend
go build
```
2. Installing Frontend NPM dependencies
```bash
cd src/frontend
npm install
```
### Compiling dockerfiles
1. Building Backend Docker:
```bash
cd src/backend
docker build -t web-brewery-element-finder-backend:latest .
```
2. Building Frontend Docker
```bash
cd src/frontend
docker build -t web-brewery-element-finder-frontend:latest .
```

## Running the Program
### Running Raw Source Files
1. Run go API  
```bash
cd src/backend
go run main.go
```
2. Run local web app
```bash
cd src/frontend
npm run dev
```
3. Open local web app  
Open a web browser and enter [http://localhost:3000/](http://localhost:3000/)
### Running Docker
1. Running Backend Docker
```bash
cd src/backend
docker run --name element-finder-backend -p 8080:8080  -d web-brewery-element-finder-backend:latest
docker start element-finder-backend 
```
2. Running Backend Docker
```bash
cd src/frontend
docker run --name element-finder-frontend -p 3000:3000 -d web-brewery-element-finder-frontend:latest
docker start element-finder-frontend
```
3. Open local web app  
Open a web browser and enter [http://localhost:3000/](http://localhost:3000/)
4. Closing Dockers
```bash
cd src/frontend
docker stop element-finder-frontend
cd src/backend
docker stop element-finder-backend 
```  
## About The Creators
<table>
  <tr>
    <th>Nama Lengkap</th>
    <th>NIM</th>
    <th>Kelas</th>
  </tr>
  <tr>
    <td>Muhammad Aufa Farabi</td>
    <td>13523023</td>
    <th>K01</th>
  </tr>
  <tr>
    <td>Darrel Adinarya Sunanda</td>
    <td>13523061</td>
    <th>K01</th>
  </tr>
  <tr>
    <td>Aryo Wisanggeni</td>
    <td>13523100</td>
    <th>K02</th>
  </tr>
</table>
