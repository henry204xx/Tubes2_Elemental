"use client";

import { useState, useEffect, useMemo } from "react";
import { Search, X } from "lucide-react";
import _ from "lodash";

export default function Searchbar({ onResults, elementsData, searchState }) {
  const [query, setQuery] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [mounted, setMounted] = useState(false);

  function isNumeric(str) {
    return !isNaN(str) && !isNaN(parseFloat(str));
  }

  function searchByName(query, elementsData) {
    if (!query) return {};
    
    return Object.fromEntries(
      Object.entries(elementsData)
        .filter(([name, _]) => 
          name.trim().toLowerCase().includes(query.toLowerCase()))
    );
  }

  function searchByTier(query, elementsData) {
    if (query === null) return {};
    
    return Object.fromEntries(
      Object.entries(elementsData)
        .filter(([_, data]) => data.tier === query)
    );
  }

  function mergeResults(nameResults, tierResults) {
    if (_.isEmpty(nameResults) || _.isEmpty(tierResults)) {
      return { ...nameResults, ...tierResults };
    }
    
    const result = {};
    Object.keys(nameResults).forEach(name => {
      if (tierResults[name]) {
        result[name] = nameResults[name];
      }
    });
    return result;
  }

  const searchData = (searchQuery) => {
    setIsLoading(true);
    
    try {
      const queries = searchQuery.trim().split(/\s+/);
      let nameQuery = '';
      let tierQuery = null;

      queries.forEach(q => {
        if (isNumeric(q)) {
          tierQuery = parseInt(q);
        } else {
          nameQuery += `${q} `;
        }
      });
      nameQuery = nameQuery.trim();

      let results = {};
      const nameResults = searchByName(nameQuery, elementsData);
      const tierResults = searchByTier(tierQuery, elementsData);

      if (nameQuery && tierQuery !== null) {
        results = mergeResults(nameResults, tierResults);
      } else {
        results = { ...nameResults, ...tierResults };
      }

      if(_.isEmpty(results) && query.length >0){
        onResults(results);
        searchState(true);
      }
      else if (!_.isEmpty(results)) {
        onResults(results);
        searchState(true);
      } 
      else{
        onResults(elementsData);
        searchState(false);
      }
    } catch (error) {
      console.error("Search error:", error);
      onResults({});
    } finally {
      setIsLoading(false);
    }
  };

  const debouncedSearch = useMemo(() =>
    _.debounce((searchQuery) => {
      searchData(searchQuery);
    }, 300), [elementsData]
  );

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    if (query.trim() !== "") {
      debouncedSearch(query);
    } else {
      onResults(elementsData);
      searchState(false);
    }

    return () => {
      debouncedSearch.cancel();
    };
  }, [query, debouncedSearch, onResults, elementsData, searchState]);

  return (
    <div className="w-full max-w-md">
      <div className="relative">
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search className="h-5 w-5 text-gray-400" />
        </div>
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Cari..."
          className="pl-10 pr-10 py-2 w-full border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        {query && (
          <button
            onClick={() => setQuery("")}
            className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
            aria-label="Clear search"
          >
            <X className="h-4 w-4" />
          </button>
        )}
      </div>
    </div>
  );
}