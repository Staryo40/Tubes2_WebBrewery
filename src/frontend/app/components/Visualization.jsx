import React, { useState } from 'react'
import { useEffect } from "react";
import elementsData from '@/public/elements.json';


const Visualization = () => {
    const [selectedElement, setSelectedElement] = useState(null);
    const [searchQuery, setSearchQuery] = useState("");
    const [filteredElements, setFilteredElements] = useState(elementsData);

    const [algorithmMode, setAlgorithmMode] = useState("BFS");
    const [isMultipleRecipe, setIsMultipleRecipe] = useState(false);
    const [isBidirectional, setIsBidirectional] = useState(false);
    const [recipeCount, setRecipeCount] = useState(1);

    const [loading, setLoading] = useState(false);
    const [result, setResult] = useState(null);
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
            setError(null);
            setResult(null);

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
                const errorData = await response.json();
                setError(errorData.Error || `HTTP error status: ${response.status}`);
                return;
              }
              const data = await response.json();
              console.log('Data dari backend:', data);
              setResult(data);
              console.log('Backend berhasil mengirim hasil resep');
            } catch (err) {
              setError('Terjadi kesalahan saat menghubungi sever.');
              console.error('Error:', err);
            } finally {
              setLoading(false);
            }


        }
  return (
    <div className="flex h-screen ">
        {/* sidebar */}
        <div className="w-1/5 flex flex-col h-screen bg-[#350535]">
          <div className="h-20 p-4 sticky top-0 z-10">
            <input
            type="text"
            placeholder="Search element"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)
              
            }
            className="w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
            />
          </div>
            <div className="flex-1 overflow-y-auto">
              <div className="space-y-2 p-4">
                {filteredElements.map((element) => (
                  <button
                      key={element.name}
                      onClick={() => {
                        setSelectedElement(element);
                        setResult(null);
                        setError(null);
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
        {/* Visualization (masih dummy) */}
        <div className="flex-1 flex-col px-6 py-4">
          {/* control bar */}
            <div className="flex justify-between items-center p-2 pb-4 border-b">
              <button
                onClick={handleStartVisualization}
                className="px-4 py-2 bg-[#5e3e5e] text-white rounded hover:bg-[#946c94] transition"
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
                onClick={() => setAlgorithmMode("BFS")}
                className={`px-4 py-2 rounded-md transition ${algorithmMode === "BFS" ? 'bg-gray-500 text-white shadow-sm' : 'text-gray-500 hover:bg-gray-50'}`}
                >
                BFS
                </button>
                <button
                onClick={() => setAlgorithmMode("DFS")}
                className={`px-4 py-2 rounded-md transition ${algorithmMode === "DFS" ? 'bg-gray-500 text-white shadow-sm' : 'text-gray-500 hover:bg-gray-50'}`}
                >
                DFS
                </button>
              </div>
            </div>
            {/* topbar tombol dfs bfs */}
            <h1 className="text-2xl font-bold text-orange-600 mb-4 mt-10">Visualization Area</h1>
            {selectedElement ? (
          <div>
            <h2 className="text-xl font-semibold mb-2">{selectedElement.name}</h2>
            <img src={selectedElement.image} alt={selectedElement.name} className="w-16 h-16 mb-2" />
            {/* Nanti diisi komponen visualisasi */}
          </div>
        ) : (
          <p className="text-gray-500">Silakan select elemen yang ingin divisualisasikan resepnya.</p>
        )}
          {loading && <p>Loading results...</p>}
          {error && <p className="text-red-500">Error: {error}</p>}
          {result && (
                    <div>
                        <h2 className="text-xl font-semibold mb-2">Recipe for {selectedElement?.name}</h2>
                        {/* Tampilkan hasil visualisasi di sini berdasarkan `visualizationResult` */}
                        <p className="mb-2 text-white">Jumlah Resep Ditemukan: {result.count}</p> 
                        <p className="mb-2 text-white">waktu: {result.elapsedTime}</p> 
                      
                    </div>
                )}
        </div>
    </div>
  )
}

export default Visualization