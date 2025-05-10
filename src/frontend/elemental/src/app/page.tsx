"use client"

import Image from "next/image";
import { useState } from "react";
import Searchbar from "./_components/searchbar.js";
import MultiValueToggle from "./_components/Checkbox.js";
import AlgorithmToggle from "./_components/algorithmTogle.js";
import Card from "./_components/card.js";
import elementLibrary from "../_dataImage/elementLibrary.json";




export default function Home() {
  const [elemToDisplay, setElemToDisplay] = useState(elementLibrary);
  const [isSearching, setIsSearching] = useState(false);
  const [isMultivalue, setMultiValue] = useState(false);
  const [graphAlgorithm, setAlgorithm] = useState("BFS");


  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen pt-12 pb-20 gap-16 px-4 sm:px-8 font-[family-name:var(--font-geist-sans)]">
      <div className="flex items-center justify-between flex-wrap w-full max-w-[100vw] px-8 overflow-x-hidden">
        <span className="text-xl font-bold text-left">ELEMENTAL</span>
        <Searchbar onResults={setElemToDisplay} elementsData={elementLibrary} searchState={setIsSearching}/>
        <MultiValueToggle onChange={setMultiValue} />
        <AlgorithmToggle onAlgorithmChange={setAlgorithm} />
      </div>

      
      <div className="flex flex-wrap gap-6 justify-center">
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
            />
          ))
        )}
      </div>
      
    </div>
  );
}
