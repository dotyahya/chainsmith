import React from 'react';

function Transaction() {
  return (
    <div>
      <h2 className="text-2xl font-bold mb-4 text-blue-400">Recent Transactions</h2>
      <ul className="space-y-4">
        <li className="bg-gray-800 p-4 rounded shadow hover:shadow-lg transition-shadow">
          <p className="text-gray-300">From: 0x123...abc</p>
          <p className="text-gray-300">To: 0x456...def</p>
          <p className="text-gray-300">Amount: 5 ETH</p>
        </li>
        <li className="bg-gray-800 p-4 rounded shadow hover:shadow-lg transition-shadow">
          <p className="text-gray-300">From: 0x789...ghi</p>
          <p className="text-gray-300">To: 0x012...jkl</p>
          <p className="text-gray-300">Amount: 2.5 ETH</p>
        </li>
        {/* Add more transactions */}
      </ul>
    </div>
  );
}

export default Transaction;
