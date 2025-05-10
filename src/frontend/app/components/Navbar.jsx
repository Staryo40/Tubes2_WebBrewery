'use client'
import React from 'react'

const Navbar = ({ onNavigate, active }) => {
  return (
    <nav className="top-0  h-auto pb-4 w-full shadow-md bg-[#380028] py-4 ">
      <div className="px-4 mx-auto flex justify-between items-center">
        <div className="mx-5 flex flex-col">
          <h1 className="text-amber-300 text-xl font-bold">Little Alchemy 2</h1>
          <p className="text-white text-sm">Recipe finder</p>
        </div>

        <div className="flex space-x-6 mr-5">
          <button
            onClick={() => onNavigate("elements")}
            className={`font-medium transition-colors ${
              active === "elements" ? "text-amber-300" : "text-white hover:text-amber-300"
            }`}
          >
            Elements
          </button>
          <button
            onClick={() => onNavigate("visualization")}
            className={`font-medium transition-colors ${
              active === "visualization" ? "text-amber-300" : "text-white hover:text-amber-300"
            }`}
          >
            Recipe Visualization
          </button>
        </div>
      </div>
    </nav>
  )
}

export default Navbar
