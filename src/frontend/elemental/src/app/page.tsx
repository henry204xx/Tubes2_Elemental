"use client";

import { useState } from "react";
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
  const [graphAlgorithm, setAlgorithm] = useState("BFS");
  const [showTreeModal, setShowTreeModal] = useState(false);
  const [showInputModal, setShowInputModal] = useState(false); 
  const [inputValue, setInputValue] = useState(1); 
  const [inputError, setInputError] = useState("");


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
    setShowInputModal(false); 
    setInputValue(1); 
    setInputError(""); 
    setIsMultiValue(false); 
  };


  const handleReceiptClick = (val: boolean) => {
      setShowTreeModal(true);
  };

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen pt-12 pb-20 gap-16 px-4 sm:px-8 font-[family-name:var(--font-geist-sans)]">
      <div className="flex items-center justify-between flex-wrap w-full max-w-[100vw] px-8 overflow-x-hidden">
        <span className="text-xl font-bold text-left">ELEMENTAL</span>
        <Searchbar onResults={setElemToDisplay} elementsData={elementLibrary} searchState={setIsSearching} />
        <MultiValueToggle onChange={handleMultivalueChange} checked={isMultivalue} />
        <AlgorithmToggle onAlgorithmChange={setAlgorithm} />
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
              onClickedButt={handleReceiptClick}
            />
          ))
        )}
      </div>

      <MultiValueInputModal 
        isOpen={showInputModal} 
        onClose={handleCancel} 
        inputValue={inputValue} 
        setInputValue={setInputValue} 
        onSubmit={handleSubmitValue} 
        inputError={inputError} 
      />

      <TreeModal
        isOpen={showTreeModal}
        onClose={() => {
          setShowTreeModal(false);
        }}
      />

    </div>
  );
}
