'use client'
import React, { useState } from 'react';
import { useEffect } from "react";
import elementsData from '@/public/elements.json';

const Elements = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedElement, setSelectedElement] = useState(null);
  const [isPopUpOpen, setIsPopUpOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [filteredElements, setFilteredElements] = useState(elementsData);

  useEffect(() => {
    const filtered = elementsData.filter(el =>
      el.name.toLowerCase().includes(searchQuery.toLowerCase())
    );
    setFilteredElements(filtered);
    setCurrentPage(1); 
  }, [searchQuery]);

  const elementsPerPage = 18
  const totalPages = Math.ceil(filteredElements.length / elementsPerPage);
  const currentElements = filteredElements.slice(
    (currentPage - 1) * elementsPerPage,
    currentPage * elementsPerPage
  );

  const handleCardClick = (element) => {
    setSelectedElement(element);
    setIsPopUpOpen(true);
  };

  const recipesPopUp = (recipes) => {
    return recipes.map((recipe, index) => (
      <div key={index} className="mb-2">
        {recipe[0] === 'Available from the start' ? (
          <span className="bg-gray-300 px-2 py-1 rounded text-black">Available from the start</span>
        ) : (
          <div className="flex items-center space-x-2">
            {recipe.map((item, i) => (
              <React.Fragment key={i}>
                <span className="bg-gray-300 px-2 py-1 rounded text-black">{item}</span>
                {i < recipe.length - 1 && <span className = "text-black">+</span>}
              </React.Fragment>
            ))}
          </div>
        )}
      </div>
    ));
  };

  return (
    <div className="mx-4 px-4 py-6 justify-center">
      <div className="mb-6">
        <input
          type="text"
          placeholder="Search element"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full p-2 border bg-[#3f383f] border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-orange-500"
          />
      </div>
 
      {/* Grid Element Cards */}
      <div className="grid grid-cols-6 gap-4">
        {currentElements.map((element) => (
          <div 
            key={element.name}
            onClick={() => handleCardClick(element)}
            className="bg-white rounded-lg shadow-md p-4 flex flex-col items-center cursor-pointer hover:shadow-lg transition-all hover:scale-105"
          >
            <img 
              src={element.image} 
              alt={element.name}
              className="w-16 h-16 object-contain mb-2"
            />
            <p className="text-center font-medium text-gray-800">{element.name}</p>
            <p className="text-xs text-gray-500">{element.category}</p>
          </div>
        ))}
      </div>

      {/* Popup */}
      {isPopUpOpen && selectedElement && (
        <div className="fixed inset-0 backdrop-blur-sm bg-black/30 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold text-orange-500">{selectedElement.name}</h2>
              <button 
                onClick={() => setIsPopUpOpen(false)}
                className="text-gray-500 hover:text-gray-700"
              >
                ✕
              </button>
            </div>
                 
            <h3 className="font-semibold mb-2 text-black">Recipes:</h3>
            <div className="space-y-2 max-h-[400px] overflow-y-auto pr-2">
              {recipesPopUp(selectedElement.recipes)}
            </div>
          </div>
        </div>
      )}

      {/* Pagination */}
        <div className="flex justify-center items-center mt-8 space-x-4">
            <button
                onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                disabled={currentPage === 1}
                className={`px-4 py-2 rounded ${
                    currentPage === 1 
                        ? 'bg-gray-500 text-white' 
                        : 'bg-white text-black border border-gray-300 hover:bg-gray-100'
                    }`}
                >
                Prev
            </button>
  
            <span className="px-4 py-2">
                {currentPage} / {totalPages}
            </span>
  
            <button
                onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                disabled={currentPage === totalPages}
                className={`px-4 py-2 rounded ${
                    currentPage === totalPages
                        ? 'bg-gray-500 text-white' 
                        : 'bg-white text-black border border-gray-300 hover:bg-gray-100'
                }`}
                >
                Next
            </button>
        </div>      
    </div>
  );
};

export default Elements;