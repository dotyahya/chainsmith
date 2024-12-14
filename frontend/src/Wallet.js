import React from 'react';

function Wallet() {
  return (
    <div className="text-center">
      <h2 className="text-2xl font-bold mb-4">Your Wallet</h2>
      <p className="text-gray-300">Balance: 10 ETH</p>
      <button className="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded mt-4">Send ETH</button>
    </div>
  );
}

export default Wallet;
