"use client"

import React from "react";
import TreeElement from "./TreeElement.js";

const TreeModal = ({ isOpen, onClose,  target, treeRaw, countNode, countSolution, programTime}) => {
  if (!isOpen) return null;
  if(!target){
    console.log("INI DI MODAL KOSONG DIE");
  }
  else{
    console.log("INI DI MODAL DIE");
    console.log(target);
  }

  return (
    <div 
      className="fixed inset-0 flex items-center justify-center z-50"
      style={{ backgroundColor: 'rgba(209, 213, 219, 0.7)' }} 
    >
      <div className="bg-white rounded-lg shadow-lg relative ">
    
        <button 
          onClick={onClose}
          className="absolute -top-5 -right-5 bg-red-600 text-white w-6 h-6 rounded-full flex items-center justify-center focus:outline-none hover:bg-red-700"
        >
          ✕
        </button>
        
        <TreeElement treeRawData={treeRaw} rootName={target} nodeCount={countNode} 
        solutionCount={countSolution} time={programTime}/>
      </div>
    </div>
  );
};

export default TreeModal;