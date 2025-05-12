import React, { useState } from 'react'
import { useEffect } from "react";
import Tree from 'react-d3-tree';
import elementsData from '@/public/elements.json';


const Visualization = () => {
    const [selectedElement, setSelectedElement] = useState(null);
    const [searchQuery, setSearchQuery] = useState("");
    const [filteredElements, setFilteredElements] = useState(elementsData);

    const [algorithmMode, setAlgorithmMode] = useState("BFS");
    const [isMultipleRecipe, setIsMultipleRecipe] = useState(false);
    const [isBidirectional, setIsBidirectional] = useState(false);
    const [recipeCount, setRecipeCount] = useState(1);
    const [currentPage, setCurrentPage] = useState(0);
    const recipesPerPage = 1;

    const [loading, setLoading] = useState(false);
    const [result, setResult] = useState(null);
    const [isStartClicked, setIsStartClicked] = useState(false);
    const [error, setError] = useState(null); 
    
    const handleMutipleCheckboxChange = (e) => {
      const isChecked = e.target.checked;
      setIsMultipleRecipe(isChecked);

      if (!isChecked) {
        setRecipeCount(1);
      } else {
        setRecipeCount(2);
      }
    }
    // useEffect buat searchbar
    useEffect(() => {
      const filtered = elementsData.filter(e =>
      e.name.toLowerCase().includes(searchQuery.toLowerCase())
      );
      setFilteredElements(filtered);
    }, [searchQuery]);

    const handleStartVisualization = async () => {
        if (!selectedElement) {
            alert('Elemen belum diseleksi');
            return;
        }

            setLoading(true);
            setIsStartClicked(true);
            setError(null);
            setResult(null);
            console.log(algorithmMode);
            console.log('tes');

            const payload = {
              Target: selectedElement.name,
              Method: algorithmMode,
              PathNumber: recipeCount,
              Bidirectional: isBidirectional,
            };

            try {
              const response = await fetch('http://localhost:8080/api/recipe', {
                method: 'POST',
                headers: {
                  'Content-type': 'application/json',
                },
                body: JSON.stringify(payload),
              });

              if (!response.ok) {
                const textData = await response.text();
                if (textData.includes('No paths found for target')) {
                    setResult(null);
                    // setError(`Tidak ada resep valid ditemukan untuk ${selectedElement?.name}.`); 
                } else {
                    setError(textData || `HTTP error status: ${response.status}`);
                }
                return;
            }
          
            const contentType = response.headers.get('Content-Type');
            if (contentType && contentType.includes('application/json')) {
                const data = await response.json();
                console.log('Data dari backend:', data);
                setResult(data);
                console.log('Backend berhasil mengirim hasil resep');
            } else {
                const textData = await response.text();
                setError(`Respon dari server bukan JSON: ${textData}`);
                console.error('Respon bukan JSON:', textData);
            }
            } catch (err) {
              setError('Terjadi kesalahan saat menghubungi server.');
              console.error('Error:', err);
            } finally {
              setLoading(false);
            }

        }
        const convertPathToTreeData = (path) => {

          // Kasus tier 0
          const targetNode = path[path.length -1];
          if (!targetNode.Ingredient1 && !targetNode.Ingredient2) {
            return [{id: targetNode.Name, name: targetNode.Name }]
          }
          const nodeMap = new Map();
          const createdNodes = new Map();
                
          const createNode = (name, isClone = false) => {
            const nodeId = `${name}-${isClone ? 'clone-' + Date.now() : 'original'}`;
            const newNode = { id: nodeId, name, children: [], isClone };
            createdNodes.set(nodeId, newNode);
            nodeMap.set(name, newNode);
            return newNode;
          };
        
          const getNode = (name) => {
            if (!nodeMap.has(name)) {
              createNode(name);
            }
            return createdNodes.get(nodeMap.get(name).id);
          };
        
          const buildingNodes = new Set(); 
        
          const buildTree = (step, currentIndex) => {
            const currentNode = getNode(step.Name);
          
            if (buildingNodes.has(currentNode.id)) {
              return currentNode; // Hindari loop
            }
            buildingNodes.add(currentNode.id);
          
            const ingredient1Name = step.Ingredient1;
            const ingredient2Name = step.Ingredient2;
          
            const ingredient1Node = getNode(ingredient1Name);
            const ingredient2Node = getNode(ingredient2Name);
          
            if (!currentNode.children.some(child => child.id === ingredient1Node.id)) {
              currentNode.children.push(ingredient1Node);
            }
        {
              currentNode.children.push(ingredient2Node);
            }
          
            // Cari langkah sebelumnya di path untuk bahan-bahan ini (sebelum currentIndex)
            for (let i = currentIndex - 1; i >= 0; i--) {
              if (path[i].Name === ingredient1Name && !ingredient1Node.children.length) {
                buildTree(path[i], i);
              }
              if (path[i].Name === ingredient2Name && !ingredient2Node.children.length) {
                buildTree(path[i], i);
              }
            }
          
            buildingNodes.delete(currentNode.id); // Hapus dari set setelah selesai membangun cabangnya
            return currentNode;
          };
        
          const treeData = [];
          if (path.length > 0) {
            treeData.push(buildTree(path[path.length - 1], path.length - 1)); // Mulai dari elemen target dengan indeksnya
          }
        
          return treeData;
        };




  return (
    <div className="flex min-h-[110vh] ">
        {/* sidebar */}
        <div className="w-1/5 flex flex-col h-[110vh] bg-[#350535]">
          <div className="h-20 p-4 top-0 z-10">
            <input
            type="text"
            placeholder="Search element"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)
              
            }
            className="w-full p-2 bg-[#3f383f] border border-gray-300 rounded focus:outline-none ring-white focus:ring-2 focus:ring-orange-500"
            />
          </div>
            <div className="flex-1 overflow-y-auto">
              <div className="space-y-2 p-4">
                {filteredElements.map((element) => (
                  <button
                      key={element.name}
                      onClick={() => {
                        setSelectedElement(element);
                        setIsStartClicked(false);
                        setResult(null);
                        setError(null);
                        setCurrentPage(0);
                      }}
                      className="flex items-center w-full rounded space-x-3 hover:bg-orange-100 transition"
                      >
                      <img 
                          src={element.image}
                          alt={element.name}
                          className="w-8 h-8 object-contain"/>
                      <span className="text-md text-gray">{element.name}</span>
                  </button>
                  ))}
              </div>
            </div>

        </div>
        {/* Visualization */}
        <div className="flex flex-col flex-1 px-6 py-4 overflow-hidden h-[110vh]">
          {/* control bar */}
            <div className="flex justify-between items-center p-2 pb-4 border-b">
             <button
               onClick={handleStartVisualization}
               className="px-4 py-2 bg-[#faa620] text-black text-bold shadow-md rounded hover:bg-[#aaa45c] transition-all hover:scale-105 active:opacity-75"
               disabled={loading}
             >
               {loading ? 'Loading...' : 'Start'}
             </button>
              <div>

              </div>
              <div className="flex flex-col items-start space-y-1">
                 <div className="grid grid-cols-2 gap-2 w-full">
                  <div className="flex items-center space-x-2">
                    <input 
                      type="checkbox"
                      checked={isMultipleRecipe}
                      onChange={handleMutipleCheckboxChange}
                      className="h-4 w-4 text-white rounded focus:ring-orange-50"
                    />
                    <span>Multiple Recipe</span>
                  </div>

                  {isMultipleRecipe && (
                    <div className="flex justify-center">
                      <input
                        type="number"
                        min="2"
                        max="100"
                        value={recipeCount}
                        onChange={(e) => {
                          const value = parseInt(e.target.value);
                          setRecipeCount(isNaN(value) ? 2 : Math.min(100, Math.max(2, value)));
                        }}
                        className="w-14 p-0.5 text-center border rounded"
                      />
                    </div>
                  )}
                </div>
                <div className="flex items-center space-x-4">
                  <label className="flex items-center space-x-2">
                    <input 
                        type="checkbox"
                        checked={isBidirectional}
                        onChange={(e) => setIsBidirectional(e.target.checked)}
                        className="h-4 w-4 text-white rounded focus:ring-orange-50"
                        />
                    <span>Bidirectional Mode</span>
                  </label>
                </div>
              </div>            

              <div className="flex items-center bg-gray-100 rounded-lg p-1">
                <button
                onClick={() => {
                  setAlgorithmMode("BFS");
                  setCurrentPage(0);
                }}
                className={`px-4 py-2 rounded-md transition ${algorithmMode === "BFS" ? 'bg-gray-500 text-white shadow-sm' : 'text-gray-500 hover:bg-gray-50'}`}
                >
                BFS
                </button>
                <button
                onClick={() => {
                  setAlgorithmMode("DFS");
                  setCurrentPage(0);
                }}
                className={`px-4 py-2 rounded-md transition ${algorithmMode === "DFS" ? 'bg-gray-500 text-white shadow-sm' : 'text-gray-500 hover:bg-gray-50'}`}
                >
                DFS
                </button>
              </div>
            </div>
            {/* topbar tombol dfs bfs */}
            <div className="mt-4"></div>
             {selectedElement ? (
              <div className="flex items-center space-x-3 px-2">
                <h2 className="text-2xl font-semibold text-amber-600 mb-2">{selectedElement.name}</h2>
                <img src={selectedElement.image} alt={selectedElement.name} className="w-12 h-12 mb-1" />
              </div>
            ) : (
              <p className="text-gray-400 mt-2 px-2">Silakan select elemen yang ingin divisualisasikan resepnya.</p>
            )}
          {loading && <p>Loading results...</p>}
          {error && <p className="text-red-500">Error: {error}</p>}
          {(!result && !loading && !isStartClicked) && (
            <p className="text-gray-300 text-xl mt-4 px-2">Belum ada visualisasi</p>
          )}
          {result === null && isStartClicked && !loading? (
            <p className="text-gray-500 py-2">Tidak ada resep valid ditemukan untuk {selectedElement?.name}.</p>
          ) : (
            result?.paths?.length > 0 ? (
              <div className="flex-1 p-2 px-2 space-y-4">
                <p className="mb-1 text-white text-xs">Jumlah Resep Ditemukan: {result.count}</p>
                <p className="mb-4 text-white text-xs">Waktu: {result.elapsedTime} ms</p>
            
                {result?.paths?.length > 1 && isMultipleRecipe && (
                  <div className="flex justify-center space-x-4 mb-4 items-center">
                    <button
                      onClick={() => setCurrentPage(currentPage - 1)}
                      disabled={currentPage === 0}
                      className={`px-3 py-1 rounded ${
                        currentPage === 0
                          ? 'bg-gray-500 text-white'
                          : 'bg-white text-black border border-gray-300 hover:bg-gray-300'
                      }`}
                    >
                      Prev
                    </button>
                    <span className="text-white">{currentPage + 1} / {result?.paths?.length}</span>
                    <button
                      onClick={() => setCurrentPage(currentPage + 1)}
                      disabled={currentPage === (result?.paths?.length - 1)}
                      className={`px-3 py-1 rounded ${
                        currentPage === (result?.paths?.length - 1)
                          ? 'bg-gray-500 text-white'
                          : 'bg-white text-black border border-gray-300 hover:bg-gray-300'
                      }`}
                    >
                      Next
                    </button>
                  </div>
                )}

                <div className="space-y-4 flex-1 ">
                  {(() => {
                    if (!result?.paths) return [];
                    const indexOfLastRecipe = (currentPage + 1) * recipesPerPage;
                    const indexOfFirstRecipe = indexOfLastRecipe - recipesPerPage;
                    const currentRecipes = result.paths.slice(indexOfFirstRecipe, indexOfLastRecipe);
                    return currentRecipes.map((recipeObject, index) => {
                      const treeData = convertPathToTreeData(recipeObject.path);
                      return treeData && (
                        <div key={index} className="border p-4 rounded-md flex flex-col flex-1 overflow-hidden">
                          <div className="flex justify-between">
                            <h3 className="font-semibold">Recipe {indexOfFirstRecipe + index + 1}</h3>
                            <p className="text-xs text-gray-300 mt-1">Gunakan pinch untuk zoom in/out dan seret untuk bergerak</p>
                          </div>
                          <p className="text-xs">Node Count: {recipeObject.nodeCount}</p>
                          <div className="overflow-auto flex-1 ">
                            <Tree
                              data={treeData}
                              orientation="bottom-top"
                              pathClassFunc={() => 'custom-link'}
                              panOnDrag={false}
                              zoomable={true}
                              translate={{ x: 420, y: 40 }}
                              depthFactor={85}
                              separation={{ siblings: 1.8, nonSiblings: 1.8 }}
                              nodeSize={{ x: 100, y: 60 }}
                              scaleExtent={{ min: 0.2, max: 0.8 }}
                              renderCustomNodeElement={({ nodeDatum }) => (
                                <g>
                                  <rect
                                    width={140}
                                    height={30}
                                    x={-70}
                                    y={-12}
                                    fill={nodeDatum.children && nodeDatum.children.length > 0 ? 'gray ' : 'white'}
                                    stroke="#ccc"
                                    style={{ cursor: 'pointer' }}
                                  />
                                  <text
                                    x={0}
                                    y={5}
                                    fill="white"
                                    textAnchor="middle"
                                    alignmentBaseline="middle"
                                    fontSize={12}
                                    letterSpacing="2px"
                                  >
                                    {nodeDatum.name}
                                  </text>
                                </g>
                              )}
                            />
                          </div>
                        </div>
                      );
                    });
                  })()}
                </div>
              </div>
            ) : null
          )}
        </div>
    </div>
  )
}

export default Visualization