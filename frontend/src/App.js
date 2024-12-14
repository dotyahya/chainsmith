import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Home from './Home';
import BlockExplorer from './BlockExplorer';
import Wallet from './Wallet';
import Transaction from './Transaction';

function App() {
  return (
    <Router>
      <div className="bg-gray-900 text-white min-h-screen">
        <nav className="bg-gray-800 p-4 flex justify-between items-center shadow-lg">
          <div className="text-2xl font-bold">Chainsmith</div>
          <ul className="flex space-x-4">
            <li><Link to="/" className="hover:text-gray-300">Home</Link></li>
            <li><Link to="/block-explorer" className="hover:text-gray-300">Block Explorer</Link></li>
            <li><Link to="/wallet" className="hover:text-gray-300">Wallet</Link></li>
            <li><Link to="/transaction" className="hover:text-gray-300">Transaction</Link></li>
          </ul>
        </nav>
        <div className="container mx-auto py-10">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/block-explorer" element={<BlockExplorer />} />
            <Route path="/wallet" element={<Wallet />} />
            <Route path="/transaction" element={<Transaction />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;
