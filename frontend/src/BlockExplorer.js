import React, { useState, useEffect } from 'react';

function BlockExplorer() {
  const [blocks, setBlocks] = useState([]);

  useEffect(() => {
    fetch('http://localhost:3001/api/blocks')
      .then(response => response.json())
      .then(data => setBlocks(data))
      .catch(err => console.error('Error fetching blocks:', err));
  }, []);

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Block Explorer</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {blocks.map(block => (
          <div key={block.id} className="bg-gray-800 p-4 rounded shadow hover:shadow-lg transition-shadow">
            <h3 className="font-bold text-lg">Block #{block.id}</h3>
            <p className="text-sm text-gray-400">Timestamp: {block.timestamp}</p>
            <p className="text-sm text-gray-400">Hash: {block.hash}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

export default BlockExplorer;
