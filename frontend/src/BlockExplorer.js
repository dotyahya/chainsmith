import React, { useState, useEffect } from 'react';

function BlockExplorer() {
  const [blocks, setBlocks] = useState([]);

  useEffect(() => {
    // Dummy backend API call for blocks
    fetch('http://localhost:5000/api/blocks')  // Replace with actual endpoint
      .then((res) => res.json())
      .then((data) => setBlocks(data.blocks));
  }, []);

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <h1 className="text-3xl font-bold mb-4">Blockchain Explorer</h1>
      <div className="space-y-4">
        {blocks.map((block) => (
          <div key={block.id} className="bg-white p-4 rounded-lg shadow-md">
            <p><strong>Block ID:</strong> {block.id}</p>
            <p><strong>Transactions:</strong> {block.transactions.length}</p>
            <a href={`/block/${block.id}`} className="text-blue-500">View Block</a>
          </div>
        ))}
      </div>
    </div>
  );
}

export default BlockExplorer;
