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
  const [isLoading, setIsLoading] = useState(false);
  const [abortController, setAbortController] = useState<AbortController | null>(null);
  const [showNoSolutionModal, setShowNoSolutionModal] = useState(false);

  useEffect(() => {
    if (treeData) {
      setIsLoading(false);
      setShowTreeModal(true);  
    } else {
      setShowTreeModal(false); 
    }
  }, [treeData]);

  useEffect(() => {
    if (manySolution == -1) {
      setIsLoading(false);
      console.log("BRO MASUK MIN 1")
      resetVal
      setShowNoSolutionModal(true);
    } 
  }, [manySolution]);

  const handleMultivalueChange = (val : boolean) => {
    if (val) {
      setShowInputModal(true);
    } else {
      setInputValue(1); 
      setIsMultiValue(false);
    }
  };


  const handleSubmitValue = () => {
    if (typeof inputValue === "number" && inputValue >= 1) {
      setIsMultiValue(true); 
      setShowInputModal(false);
      setInputError("")
    } else {
      setInputError("Nilai harus lebih dari 0 atau pastikan masukan adalah angka bulat");
    }
  };

  const resetVal = ()=>{
    setTreeData(null);
    setTarget("");
    setNodeCount(0);
    setSolutionCount(0);
    setTimeCount(0);
  }


  const handleCancel = () => {
    resetVal
    setInputValue(1); 
    setInputError(""); 
    setIsMultiValue(false); 
    setShowTreeModal(false);
    setShowInputModal(false);
  };

  const handleTreeClose = () => {
    resetVal
    setShowTreeModal(false);
  };




    const handleReceiptClick = async (elementName: string, maxSolution: number, method: string) => {
      const controller = new AbortController();
      setAbortController(controller);
      setIsLoading(true);

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
        onInputEror={setInputError}
      />

      <TreeModal
        isOpen={showTreeModal}
        onClose={handleTreeClose}
        target={currTarget} 
        treeRaw={treeData}    
        countNode={nodeCount}
        countSolution={manySolution}
        programTime={timeCount}
      />

      {/*modal prosess query*/}
      {isLoading && (
      <div className="fixed inset-0 flex items-center justify-center z-50"
      style={{ backgroundColor: 'rgba(209, 213, 219, 0.7)' }}>
        <div className="bg-white p-6 rounded-xl shadow-lg text-center">
          <p className="text-lg font-semibold mb-4">Memproses Query...</p>
          <button
            onClick={() => {
            if (abortController) abortController.abort();
              setIsLoading(false);
              handleCancel();
            }}
            className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
            >
            Batalkan
          </button>
          </div>
      </div>)}

      {/*modal no solution*/}
      {showNoSolutionModal && (
      <div className="fixed inset-0 flex items-center justify-center z-50"
      style={{ backgroundColor: 'rgba(209, 213, 219, 0.7)' }}>
        <div className="bg-white p-6 rounded-xl shadow-lg text-center">
          <h2 className="text-xl font-bold mb-4 text-red-600">Tidak ada solusi ditemukan</h2>
          <p className="mb-6 text-gray-700">Silakan coba elemen lain.</p>
          <button
            onClick={() => {
              setShowNoSolutionModal(false);
              resetVal();
            }}
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
          >
            Tutup
          </button>
        </div>
      </div>
    )}

    </div>
  );
}