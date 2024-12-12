import React, { useState, useEffect } from 'react';

function Transaction() {
  const [transactions, setTransactions] = useState([]);
  const [sender, setSender] = useState('');
  const [receiver, setReceiver] = useState('');
  const [amount, setAmount] = useState('');

  useEffect(() => {
    // Dummy backend API call for transactions
    fetch('http://localhost:5000/api/transactions')  // Replace with your actual endpoint
      .then((res) => res.json())
      .then((data) => setTransactions(data.transactions));
  }, []);

  const handleSubmit = (e) => {
    e.preventDefault();
    // Dummy API for sending a transaction
    fetch('http://localhost:5000/api/send-transaction', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ sender, receiver, amount }),
    })
      .then((res) => res.json())
      .then((data) => {
        alert(data.message);
      });
  };

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <h1 className="text-3xl font-bold mb-4">Send Transaction</h1>
      <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow-md w-96">
        <input
          type="text"
          placeholder="Sender"
          value={sender}
          onChange={(e) => setSender(e.target.value)}
          className="mb-4 p-2 border rounded"
        />
        <input
          type="text"
          placeholder="Receiver"
          value={receiver}
          onChange={(e) => setReceiver(e.target.value)}
          className="mb-4 p-2 border rounded"
        />
        <input
          type="number"
          placeholder="Amount"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          className="mb-4 p-2 border rounded"
        />
        <button type="submit" className="bg-blue-500 text-white p-2 rounded">Send</button>
      </form>

      <h2 className="text-2xl font-semibold mt-6">Recent Transactions</h2>
      <div className="space-y-4">
        {transactions.map((transaction) => (
          <div key={transaction.id} className="bg-white p-4 rounded-lg shadow-md">
            <p><strong>Sender:</strong> {transaction.sender}</p>
            <p><strong>Receiver:</strong> {transaction.receiver}</p>
            <p><strong>Amount:</strong> {transaction.amount}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

export default Transaction;
