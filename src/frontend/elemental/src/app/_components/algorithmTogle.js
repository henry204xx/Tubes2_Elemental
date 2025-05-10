"use client"

import { useState } from 'react';

export default function AlgorithmToggle({ onAlgorithmChange }) {
  const [algorithm, setAlgorithm] = useState('BFS'); // 'BFS' atau 'DFS'

  const handleToggle = (selectedAlgorithm) => {
    setAlgorithm(selectedAlgorithm);
    if (onAlgorithmChange) {
      onAlgorithmChange(selectedAlgorithm);
    }
  };

  return (
    <div className="flex flex-col items-start gap-2">
      <span className="text-sm font-medium text-gray-700">Algoritma Pencarian</span>
      <div className="flex rounded-md shadow-sm" role="group">
        <button
          type="button"
          onClick={() => handleToggle('BFS')}
          className={`px-4 py-2 text-sm font-medium rounded-l-lg border border-gray-200 ${
            algorithm === 'BFS'
              ? 'bg-blue-500 text-white'
              : 'bg-white text-gray-700 hover:bg-gray-50'
          }`}
        >
          BFS
        </button>
        <button
          type="button"
          onClick={() => handleToggle('DFS')}
          className={`px-4 py-2 text-sm font-medium rounded-r-lg border border-gray-200 ${
            algorithm === 'DFS'
              ? 'bg-blue-500 text-white'
              : 'bg-white text-gray-700 hover:bg-gray-50'
          }`}
        >
          DFS
        </button>
      </div>
      <div className="text-xs text-gray-500">
        {algorithm === 'BFS' 
          ? "Breadth-First Search"
          : "Depth-First Search"}
      </div>
    </div>
  );
}