<!-- Back to Top Link-->
<a name="readme-top"></a>


<br />
<div align="center">
  <h1 align="center"> Little Alchemy 2 Solver Using BFS and DFS
</h1>

  <p align="center">
    <h4>Tugas Besar 2 IF2211 Strategi Algoritma</h4>

  </p>
</div>

<!-- CONTRIBUTOR -->
<div align="center">
  <strong>
    <h3>Made By:</h3>
    <h3>Elemental</h3>
    <table align="center">
      <tr>
        <td>NIM</td>
        <td>Nama</td>
      </tr>
      <tr>
        <td>12821046</td>
        <td>Fardhan Indrayesa</td>
      </tr>
      <tr>
        <td>13523051</td>
        <td>Ferdinand Gabe Tua Sinaga</td>
      </tr>
      <tr>
        <td>13523108</td>
        <td>Henry Filberto Shenelo</td>
      </tr>
    </table>
  </strong>
  <br>
</div>


## External Links

- [Spesifikasi](https://docs.google.com/document/d/1aQB5USxfUCBfHmYjKl2wV5WdMBzDEyojE5yxvBO3pvc/edit?tab=t.0)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ABOUT THE PROJECT -->
## About The Project

 Solving the Little Alchemy 2 game to find the recipe of an element, from basic elements to the desired element. This program uses the Go programming language to perform recipe scraping and utilizes Breadth-First Search and Depth-First Search algorithms in searching for the recipe.

DFS utilizes a stack to explore each possible path deeply before backtracking.

BFS utilizes a queue to explore all possible paths level by level, ensuring the shortest path is found first.

 The results are visualized on a website in the form of a tree.
  

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- GETTING STARTED --> 
## Getting Started 

### Prerequisites

Project dependencies

* Node  
  You can find Node here: https://nodejs.org/en/learn/getting-started/how-to-install-nodejs  

* Docker  
  You can find Docker here: https://www.docker.com/get-started/  

* Go  
  You can find Go here: https://go.dev/doc/install  

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Installation

How to install and use this project (without docker)

1. Clone the repo
   ```sh
   git clone https://github.com/henry204xx/Tubes2_Elemental
   ```
2. Open terminal at Tubes2_Elemental directory and type
    ```sh
   cd src/backend
   ```
   and run 
    ```sh
   go run main.go
   ```

3. On other terminal, type
    ```sh
   cd src/frontend/elemental
   ``` 
4. Install node dependencies
   ```sh
   npm install
   ```

5. Run program with
   ```sh
   npm run dev
   ```
6. Now, the web app is running with the frontend and backend functionality. Visit `http://localhost:3000/` to open the Little Alchemy 2 Solver Web.
<br>

How to install and use this project (With docker)

1. Clone the repo
   ```sh
   git clone https://github.com/henry204xx/Tubes2_Elemental
   ```

2. Open your terminal at the Tubes2_Elemental directory, type
   ```sh
   docker compose up --build
   ```
3. If you want to run the web app again and had done step 2, just type
    ```sh
    docker compose up
    ```
4. Now, the web app is running with the frontend and backend functionality. Visit `http://localhost:3000/` to open the Little Alchemy 2 Solver Web.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- FEATURES -->
## Features
1. Single Recipe Search with BFS or DFS
2. Multiple Recipes Search With BFS or DFS
3. Recipe visualization with tree structure
4. Optimized With Multithreading

<p align="right">(<a href="#readme-top">back to top</a>)</p>
