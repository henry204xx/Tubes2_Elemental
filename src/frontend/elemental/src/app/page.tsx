"use client";

import { useState, useEffect } from "react";
import Searchbar from "./_components/searchbar";
import MultiValueToggle from "./_components/Checkbox";
import AlgorithmToggle from "./_components/algorithmTogle";
import Card from "./_components/card";
import elementLibrary from "../_dataImage/elementLibrary.json";
import MultiValueInputModal from "./_components/multiValueInpputModal.js";
import TreeModal from "./_components/TreeModal.js";


export default function Home() {
  const [elemToDisplay, setElemToDisplay] = useState(elementLibrary);
  const [isSearching, setIsSearching] = useState(false);
  const [isMultivalue, setIsMultiValue] = useState(false);
  const [graphAlgorithm, setAlgorithm] = useState("bfs");
  const [showTreeModal, setShowTreeModal] = useState(false);
  const [showInputModal, setShowInputModal] = useState(false); 
  const [inputValue, setInputValue] = useState(1); 
  const [inputError, setInputError] = useState("");
  const [treeData, setTreeData] = useState(null);
  const [currTarget, setTarget] = useState("");
  const [nodeCount, setNodeCount] = useState(0);
  const [manySolution, setSolutionCount] = useState(0);
  const [timeCount, setTimeCount] = useState(0);

  useEffect(() => {
    if (treeData) {
      setShowTreeModal(true);  
    } else {
      setShowTreeModal(false); 
    }
  }, [treeData]);

  const handleMultivalueChange = (val : boolean) => {
    if (val) {
      setShowInputModal(true);
    } else {
      setInputValue(1); 
      setIsMultiValue(false);
    }
  };


  const handleSubmitValue = () => {
    if (inputValue >= 1) {
      setIsMultiValue(true); 
      setShowInputModal(false);
    } else {
      setInputError("Nilai harus lebih dari 0");
    }
  };


  const handleCancel = () => {
    setTreeData(null);
    setInputValue(1); 
    setTarget("");
    setInputError(""); 
    setNodeCount(0);
    setSolutionCount(0);
    setTimeCount(0);
    setIsMultiValue(false); 
    setShowTreeModal(false);
  };




    const handleReceiptClick = async (elementName: string, maxSolution: number, method: string) => {
      try {
        console.log("MASUK PROSES");
        console.log("INI DIBAWAHNYA NAMANYA");
        console.log(elementName);
        console.log("INI DIBAWAHNYA METHODNYA");
        console.log(method);
        console.log("INI TYPE MAX SOLUTION "+typeof(maxSolution));
        const response = await fetch('http://localhost:8080/api/query', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
        body: JSON.stringify({
          target: elementName,
          maxSolutions: maxSolution,
          method: method
        }),
      });


      if (!response.ok) throw new Error('Failed to fetch tree data');
      
      const data = await response.json();
      setTreeData(data.trees);
      setTarget(elementName);
      setNodeCount(data.nodeCount);
      setSolutionCount(data.numSolutions);
      setTimeCount(data.elapsedTime)
      

      console.log("Berhasil");
    } catch (error) {
      console.error("Error fetching tree data:", error);
    } 
  };

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen pt-12 pb-20 gap-16 px-4 sm:px-8 font-[family-name:var(--font-geist-sans)]">
      <div className="flex items-center justify-between flex-wrap w-full max-w-[100vw] px-8 overflow-x-hidden">
        <span className="text-xl font-bold text-left">ELEMENTAL</span>
        <Searchbar onResults={setElemToDisplay} elementsData={elementLibrary} searchState={setIsSearching} />
        <MultiValueToggle onChange={handleMultivalueChange} checked={isMultivalue} />
        <AlgorithmToggle onToggle={setAlgorithm} algorithm={graphAlgorithm} />
      </div>
      <div className="flex flex-wrap gap-6 justify-center ">
        {isSearching && Object.keys(elemToDisplay).length === 0 ? (
          <div className="text-center py-10 w-full">
            <p className="text-gray-500 text-lg">Maaf, item tidak ditemukan</p>
          </div>
        ) : (
          Object.entries(elemToDisplay).map(([name, { svgPath, tier }]) => (
            <Card 
              key={name} 
              title={name} 
              image={svgPath ? `../${svgPath}` : ''} 
              tier={tier} 
              onClickedButt={() => handleReceiptClick(name, inputValue, graphAlgorithm)}
            />
          ))
        )}
      </div>

      <MultiValueInputModal 
        isOpen={showInputModal} 
        onClose={handleCancel} 
        setInputValue={setInputValue} 
        onSubmit={handleSubmitValue} 
        inputError={inputError} 
      />

      <TreeModal
        isOpen={showTreeModal}
        onClose={()=>setShowTreeModal(false)}
        target={currTarget} 
        treeRaw={treeData}    
        countNode={nodeCount}
        countSolution={manySolution}
        programTime={timeCount}
      />

    </div>
  );
}