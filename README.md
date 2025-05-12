# Tugas Besar 2 Strategi Algoritma IF2211 2025

## Little Alchemy 2 Element Recipe Finder with BFS, DFS, and Bidirectional Algorithms
Breadth first search (BFS) is a pathfinding algorithm that focuses on exploring each alternative little by little concurrently. On the other hand, Depth First Search (DFS), also a pathfinding algorithm, focuses on choosing one alternative and going through it until it reaches the end before exploring another alternative. Performance-wise, BFS guarantees in finding the most efficient path towards the goal but is cumbersome on memory-usage, while DFS is less memory hungry and faster, but is not guaranteed to find the most efficient path.  

Both algorithms are implemented in creating this simple web application project to find an element in the game called Little Alchemy 2, which is a game where we combine elements to find other more higher tier element. The goal of the game is to find all elements that are possible in the game. There are 720 elements in total in this game without any extra add-ons and finding the recipes to complete the game can become troublesome. That is why this app helps in finding paths to get to a certain element target. Specifically for this project the algorithm used are BFS, DFS, and Bidirectional search. Users can switch between modes to see the performance difference of each algorithm. In addition, there is a feature to show multiple paths towards a single target element that is implemented with multithreading to speed up the process of searcing.

## Program Requirements
### Backend (Go)
  - Version **1.18+** (ideally the latest 1.x release)  

### Frontend (Next.js)
- **Node.js**  
  - Version **16.x** or **18.x**  
- **npm** (v8+) _or_ **Yarn** (v1.22+)

### Containerized Setup
- **Docker** & **Docker Compose**  

## Compiling the Program
Compile the program only if you want to run the program locally with the raw code
### Backend (Go)  
1. Navigate to the `backend/` folder:
```bash
cd src/backend
```
2. Download go dependencies
```bash
go mod download
```
### Frontend (Next.js)
1. Navigate to the `frontend/` folter:
```bash
cd src/frontend
```
2. Install npm packages
```bash
npm intall
```

## Running the Program
### Running locally
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