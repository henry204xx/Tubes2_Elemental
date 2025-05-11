"use client"

import React, { useEffect, useRef, useState } from "react";
import * as d3 from "d3";
import elementLibrary from "../../_dataImage/elementLibrary.json";

function getSvgPath(name) {
  return elementLibrary[name]?.svgPath || "";
}

function buildOutputMap(rules) {
  const map = {};
  for (const rule of rules) {
    if (!map[rule.output]) map[rule.output] = [];
    map[rule.output].push(rule.input);
  }
  return map;
}

function buildTree(nodeName, outputMap) {
  const combinations = outputMap[nodeName] || [];
  const children = combinations.map(inputs => ({
    name: inputs.join("+"),
    children: inputs.map(input => buildTree(input, outputMap))
  }));
  return { name: nodeName, children };
}

function splitText(text, maxChars = 9) {
  return text.length <= maxChars ? [text] : [text.substring(0, maxChars)];
}

const dataRule = [

{ "input": ["Air", "Planet"], "output": "Atmosphere" },
{ "input": ["Continent", "Continent"], "output": "Planet" },
{ "input": ["Land", "Land"], "output": "Continent" },
{ "input": ["Land", "Earth"], "output": "Continent" },
{ "input": ["Earth", "Earth"], "output": "Land" },


];

const outputMap = buildOutputMap(dataRule);
const treeDummy = buildTree("Atmosphere", outputMap);
console.log(treeDummy);

function TreeElement() {
  const svgRef = useRef();
  const containerRef = useRef();
  const [tree, setTree] = useState(null);
  const [zoom, setZoom] = useState(1);

  // useEffect(() => {
  //   fetch("/api/tree-data")
  //     .then(res => res.json())
  //     .then(dataRule => {
  //       const outputMap = buildOutputMap(dataRule);
  //       const treeData = buildTree("Atmosphere", outputMap);
  //       setTree(treeData);
  //     });
  // }, []);

  useEffect(() => { //ganti tree dummy dengan data hasil fetch nantinya
    if (!treeDummy) return;

    const width = 1200;
    const height = 600;

    const svg = d3.select(svgRef.current);
    svg.selectAll("*").remove();
    svg
      .attr("width", width)
      .attr("height", height)
      .attr("viewBox", [0, 0, width, height]);

    const zoomContainer = svg.append("g");

    const zoomBehavior = d3.zoom()
      .scaleExtent([0.1, 4])
      .on("zoom", event => {
        zoomContainer.attr("transform", event.transform);
        setZoom(event.transform.k);
      });

    svg.call(zoomBehavior);
    svg.on("dblclick.zoom", () => {
      svg.transition().duration(750).call(
        zoomBehavior.transform,
        d3.zoomIdentity
      );
    });

    const root = d3.hierarchy(treeDummy);
    const treeLayout = d3.tree().size([width - 100, height * 1.2]).separation(() => 1.5);
    treeLayout(root);

    const nodes = root.descendants();
    const xValues = nodes.map(d => d.x);
    const yValues = nodes.map(d => d.y);
    const xOffset = (width - (Math.max(...xValues) - Math.min(...xValues))) / 2 - Math.min(...xValues) - 600;
    const yOffset = (height - (Math.max(...yValues) - Math.min(...yValues))) / 2 - Math.min(...yValues) + 150;

    const g = zoomContainer.append("g").attr("transform", `translate(${xOffset},${yOffset})`);

    g.selectAll(".link")
      .data(root.links())
      .enter()
      .append("line")
      .attr("stroke", "#aaa")
      .attr("stroke-width", 1.5)
      .attr("x1", d => d.source.x)
      .attr("y1", d => d.source.y)
      .attr("x2", d => d.target.x)
      .attr("y2", d => d.target.y);

    const node = g
      .selectAll(".node")
      .data(root.descendants())
      .enter()
      .append("g")
      .attr("class", "node")
      .attr("transform", d => `translate(${d.x},${d.y})`);

    node.each(function (d) {
      const group = d3.select(this);
      const name = d.data.name;

      if (name.includes("+")) {
        const [left, right] = name.split("+");

        group.append("rect")
          .attr("x", -60)
          .attr("y", -40)
          .attr("width", 120)
          .attr("height", 60)
          .attr("rx", 10)
          .attr("fill", "#4b0082");

        group.append("image").attr("x", -50).attr("y", -30).attr("width", 30).attr("href", getSvgPath(left));
        group.selectAll(".left-text")
          .data(splitText(left))
          .enter()
          .append("text")
          .attr("x", -35).attr("y", (d, i) => 12 + (i * 12))
          .attr("text-anchor", "middle")
          .style("font-size", "9px")
          .style("fill", "white")
          .text(d => d);

        group.append("text").attr("x", -4).attr("y", -5).style("font-size", "16px").style("fill", "white").text("+");

        group.append("image").attr("x", 20).attr("y", -30).attr("width", 30).attr("href", getSvgPath(right));
        group.selectAll(".right-text")
          .data(splitText(right))
          .enter()
          .append("text")
          .attr("x", 35).attr("y", (d, i) => 12 + (i * 12))
          .attr("text-anchor", "middle")
          .style("font-size", "9px")
          .style("fill", "white")
          .text(d => d);
      } else {
        group.append("circle").attr("r", 25).attr("fill", "#4b0082").attr("stroke", "#9370db").attr("stroke-width", 2);
        group.append("image").attr("x", -15).attr("y", -15).attr("width", 30).attr("height", 30).attr("href", getSvgPath(name));
      }
    });

    svg.call(zoomBehavior.transform, d3.zoomIdentity.translate(width / 2, height / 10).scale(0.8));
  }, [tree]);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div ref={containerRef} style={{
        width: '600px',
        height: '300px',
        position: 'relative',
        overflow: 'hidden',
        touchAction: 'none',
        border: '1px solid #ccc',
        borderRadius: '8px'
      }}>
        <svg ref={svgRef} style={{ display: 'block', width: '100%', height: '100%', minHeight: '400px' }} />
        <div style={{
          position: 'absolute', top: '10px', left: '10px',
          background: 'rgba(255,255,255,0.7)', padding: '5px',
          borderRadius: '5px', fontSize: '12px', maxWidth: '200px', opacity: 0.8
        }}>
          Gunakan pinch untuk zoom in/out, dan seret untuk bergerak
        </div>
      </div>

      <div style={{
        padding: '15px',
        backgroundColor: '#f8f9fa',
        borderRadius: '8px',
        border: '1px solid #dee2e6',
        fontFamily: 'Arial, sans-serif'
      }}>
        <h3 style={{ marginTop: 0, color: '#4b0082' }}>Informasi Program</h3>
        <div style={{ display: 'flex', gap: '20px' }}>
          <div>
            <strong>Waktu Eksekusi:</strong> 
            <div style={{ 
              background: '#e9ecef', 
              padding: '5px 10px', 
              borderRadius: '4px',
              marginTop: '5px'
            }}>
              {0}
            </div>
          </div>
          <div>
            <strong>Jumlah Node:</strong>
            <div style={{ 
              background: '#e9ecef', 
              padding: '5px 10px', 
              borderRadius: '4px',
              marginTop: '5px'
            }}>
              {0}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default TreeElement;




