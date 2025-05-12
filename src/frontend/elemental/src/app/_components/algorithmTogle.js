"use client"

export default function AlgorithmToggle({ algorithm, onToggle }) {
  return (
    <div className="flex flex-col items-start gap-2">
      <span className="text-sm font-medium text-gray-700">Algoritma Pencarian</span>
      <div className="flex rounded-md shadow-sm" role="group">
        <button
          type="button"
          onClick={() => onToggle('bfs')}
          className={`px-4 py-2 text-sm font-medium rounded-l-lg border border-gray-200 ${
            algorithm === 'bfs'
              ? 'bg-blue-500 text-white'
              : 'bg-white text-gray-700 hover:bg-gray-50'
          }`}
        >
          BFS
        </button>
        <button
          type="button"
          onClick={() => onToggle('dfs')}
          className={`px-4 py-2 text-sm font-medium rounded-r-lg border border-gray-200 ${
            algorithm === 'dfs'
              ? 'bg-blue-500 text-white'
              : 'bg-white text-gray-700 hover:bg-gray-50'
          }`}
        >
          DFS
        </button>
      </div>
      <div className="text-xs text-gray-500">
        {algorithm === 'bfs' 
          ? "Breadth-First Search"
          : "Depth-First Search"}
      </div>
    </div>
  );
}