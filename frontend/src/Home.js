import React, { useState, useEffect } from 'react';

function Home() {
  const [blockCount, setBlockCount] = useState(null);
  const [transactionCount, setTransactionCount] = useState(null);

  useEffect(() => {
    // Dummy backend API call
    fetch('http://localhost:5000/api/status')  // Replace with your actual backend URL
      .then((res) => res.json())
      .then((data) => {
        setBlockCount(data.blockCount);
        setTransactionCount(data.transactionCount);
      });
  }, []);

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col items-center justify-center">
      <h1 className="text-3xl font-bold mb-4">Blockchain Dashboard</h1>
      <div className="bg-white p-6 rounded-lg shadow-lg w-96">
        <h2 className="text-xl">Blockchain Overview</h2>
        <p>Block Count: {blockCount ? blockCount : 'Loading...'}</p>
        <p>Transaction Count: {transactionCount ? transactionCount : 'Loading...'}</p>
      </div>
    </div>
  );
}

export default Home;
