import React, { useState, useEffect } from 'react';

function Wallet() {
  const [balance, setBalance] = useState(null);
  const [transactions, setTransactions] = useState([]);

  useEffect(() => {
    // Dummy backend API call for wallet balance and transactions
    fetch('http://localhost:5000/api/wallet')  // Replace with actual backend endpoint
      .then((res) => res.json())
      .then((data) => {
        setBalance(data.balance);
        setTransactions(data.transactions);
      });
  }, []);

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <h1 className="text-3xl font-bold mb-4">My Wallet</h1>
      <div className="bg-white p-6 rounded-lg shadow-md w-96">
        <h2 className="text-xl">Wallet Balance: {balance ? balance : 'Loading...'}</h2>
      </div>

      <h2 className="text-2xl font-semibold mt-6">Transaction History</h2>
      <div className="space-y-4">
        {transactions.map((transaction) => (
          <div key={transaction.id} className="bg-white p-4 rounded-lg shadow-md">
            <p><strong>Amount:</strong> {transaction.amount}</p>
            <p><strong>Date:</strong> {transaction.date}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

export default Wallet;
