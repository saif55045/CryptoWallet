import React, { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import api from '../services/api';
import Navbar from '../components/Navbar';

const Dashboard = () => {
  const { user, token } = useAuth();
  const [wallet, setWallet] = useState(null);
  const [balance, setBalance] = useState(0);
  const [transactions, setTransactions] = useState([]);
  const [blocksMined, setBlocksMined] = useState(0);
  const [zakatDue, setZakatDue] = useState(0);
  const [loading, setLoading] = useState(true);

  const fetchBlocksMined = useCallback(async () => {
    try {
      const res = await api.blockchain.getMyBlocks(token);
      setBlocksMined(res.data.blocks?.length || 0);
    } catch (err) {
      console.log('Blocks info not available');
    }
  }, [token]);

  const fetchZakatStatus = useCallback(async () => {
    try {
      const res = await api.zakat.getSummary(token);
      setZakatDue(res.data.summary?.zakatDue || 0);
    } catch (err) {
      console.log('Zakat info not available');
    }
  }, [token]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const walletRes = await api.wallet.getMyWallet();
        setWallet(walletRes.data.wallet);

        const balanceRes = await api.utxo.getMyBalance();
        // Balance response is { balance: { walletId, balance, confirmedBalance, pendingBalance, utxoCount } }
        const balanceValue = balanceRes.data.balance?.balance || 0;
        setBalance(Number(balanceValue) || 0);

        const txRes = await api.transaction.getMyTransactions();
        setTransactions(txRes.data.transactions?.slice(0, 5) || []);
      } catch (err) {
        console.error('Error loading dashboard:', err);
        setBalance(0);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
    fetchBlocksMined();
    fetchZakatStatus();
  }, [token, fetchBlocksMined, fetchZakatStatus]);

  const StatCard = ({ icon, title, value, subtitle, gradient, link }) => (
    <Link to={link} className="glass-card rounded-2xl p-6 card-hover block">
      <div className="flex items-start justify-between">
        <div>
          <p className="text-gray-600 text-sm font-medium">{title}</p>
          <p className={`text-3xl font-bold mt-2 ${gradient}`}>{value}</p>
          {subtitle && <p className="text-gray-500 text-sm mt-1">{subtitle}</p>}
        </div>
        <div className={`w-12 h-12 rounded-xl flex items-center justify-center text-2xl ${
          gradient.includes('green') ? 'bg-green-500/20' :
          gradient.includes('purple') ? 'bg-purple-500/20' :
          gradient.includes('yellow') ? 'bg-yellow-500/20' :
          'bg-indigo-500/20'
        }`}>
          {icon}
        </div>
      </div>
    </Link>
  );

  if (loading) {
    return (
      <div className="min-h-screen">
        <Navbar />
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <div className="w-16 h-16 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
            <p className="text-gray-600 mt-4">Loading your dashboard...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Welcome Section */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-800 mb-2">
            Welcome back, {user?.fullName?.split(' ')[0]} 👋
          </h1>
          <p className="text-gray-600">Here's what's happening with your wallet today.</p>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <StatCard
            icon="💰"
            title="Total Balance"
            value={`${balance.toFixed(2)}`}
            subtitle="Available coins"
            gradient="text-green-400"
            link="/wallet"
          />
          <StatCard
            icon="📊"
            title="Transactions"
            value={transactions.length}
            subtitle="Recent activity"
            gradient="text-indigo-400"
            link="/send"
          />
          <StatCard
            icon="⛏️"
            title="Blocks Mined"
            value={blocksMined}
            subtitle="Total mined"
            gradient="text-purple-400"
            link="/mining"
          />
          <StatCard
            icon="🕌"
            title="Zakat Due"
            value={zakatDue.toFixed(2)}
            subtitle={zakatDue > 0 ? 'Payment pending' : 'All clear'}
            gradient="text-yellow-400"
            link="/zakat"
          />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Wallet Card */}
          <div className="lg:col-span-2 glass-card rounded-2xl p-6">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-bold text-gray-800">Your Wallet</h2>
              <Link to="/wallet" className="text-indigo-600 hover:text-indigo-700 text-sm font-medium transition-colors">
                View Details →
              </Link>
            </div>

            {wallet ? (
              <div className="space-y-4">
                <div className="bg-gradient-to-r from-indigo-500/10 to-purple-500/10 rounded-xl p-4 border border-indigo-500/20">
                  <p className="text-gray-600 text-sm mb-1">Wallet Address</p>
                  <p className="text-gray-800 font-mono text-sm break-all">{wallet.walletId}</p>
                </div>
                
                <div className="grid grid-cols-2 gap-4">
                  <div className="bg-gray-100 rounded-xl p-4">
                    <p className="text-gray-600 text-sm mb-1">Balance</p>
                    <p className="text-2xl font-bold text-green-600">{balance.toFixed(4)}</p>
                  </div>
                  <div className="bg-gray-100 rounded-xl p-4">
                    <p className="text-gray-600 text-sm mb-1">Status</p>
                    <div className="flex items-center gap-2">
                      <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                      <span className="text-green-600 font-medium">Active</span>
                    </div>
                  </div>
                </div>
              </div>
            ) : (
              <div className="text-center py-8">
                <div className="w-16 h-16 bg-indigo-500/20 rounded-2xl flex items-center justify-center mx-auto mb-4">
                  <span className="text-3xl">💼</span>
                </div>
                <h3 className="text-gray-800 font-semibold mb-2">No Wallet Found</h3>
                <p className="text-gray-600 mb-4">Generate a wallet to start using CryptoVault</p>
                <Link to="/wallet" className="btn-primary inline-flex items-center gap-2">
                  Generate Wallet
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                  </svg>
                </Link>
              </div>
            )}
          </div>

          {/* Quick Actions */}
          <div className="glass-card rounded-2xl p-6">
            <h2 className="text-xl font-bold text-gray-800 mb-6">Quick Actions</h2>
            
            <div className="space-y-3">
              <Link to="/send" className="flex items-center gap-4 p-4 rounded-xl bg-gray-100 hover:bg-gray-200 transition-all group">
                <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                  </svg>
                </div>
                <div className="flex-1">
                  <p className="text-gray-800 font-medium">Send Money</p>
                  <p className="text-gray-500 text-sm">Transfer to another wallet</p>
                </div>
                <svg className="w-5 h-5 text-gray-400 group-hover:text-gray-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </Link>

              <Link to="/mining" className="flex items-center gap-4 p-4 rounded-xl bg-gray-100 hover:bg-gray-200 transition-all group">
                <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-purple-500 to-pink-600 flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                  </svg>
                </div>
                <div className="flex-1">
                  <p className="text-gray-800 font-medium">Mine Blocks</p>
                  <p className="text-gray-500 text-sm">Earn mining rewards</p>
                </div>
                <svg className="w-5 h-5 text-gray-400 group-hover:text-gray-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </Link>

              <Link to="/zakat" className="flex items-center gap-4 p-4 rounded-xl bg-gray-100 hover:bg-gray-200 transition-all group">
                <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-yellow-500 to-orange-600 flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                  </svg>
                </div>
                <div className="flex-1">
                  <p className="text-gray-800 font-medium">Pay Zakat</p>
                  <p className="text-gray-500 text-sm">Calculate & pay zakat</p>
                </div>
                <svg className="w-5 h-5 text-gray-400 group-hover:text-gray-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </Link>
            </div>
          </div>
        </div>

        {/* Recent Transactions */}
        <div className="mt-6 glass-card rounded-2xl p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-bold text-gray-800">Recent Transactions</h2>
            <Link to="/send" className="text-indigo-600 hover:text-indigo-700 text-sm font-medium transition-colors">
              View All →
            </Link>
          </div>

          {transactions.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="text-left text-gray-600 text-sm border-b border-gray-200">
                    <th className="pb-4 font-medium">Type</th>
                    <th className="pb-4 font-medium">Amount</th>
                    <th className="pb-4 font-medium hidden sm:table-cell">Status</th>
                    <th className="pb-4 font-medium hidden md:table-cell">Date</th>
                  </tr>
                </thead>
                <tbody>
                  {transactions.map((tx, index) => (
                    <tr key={tx.transactionId || index} className="border-b border-gray-200">
                      <td className="py-4">
                        <div className="flex items-center gap-3">
                          <div className={`w-8 h-8 rounded-lg flex items-center justify-center ${
                            tx.direction === 'sent'
                              ? 'bg-red-100 text-red-600'
                              : 'bg-green-100 text-green-600'
                          }`}>
                            {tx.direction === 'sent' ? '↑' : '↓'}
                          </div>
                          <span className="text-gray-800 font-medium">
                            {tx.direction === 'sent' ? 'Sent' : 'Received'}
                          </span>
                        </div>
                      </td>
                      <td className={`py-4 font-semibold ${
                        tx.direction === 'sent' ? 'text-red-600' : 'text-green-600'
                      }`}>
                        {tx.direction === 'sent' ? '-' : '+'}{tx.amount?.toFixed(2) || '0.00'}
                      </td>
                      <td className="py-4 hidden sm:table-cell">
                        <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                          tx.status === 'confirmed'
                            ? 'bg-green-100 text-green-600'
                            : 'bg-yellow-100 text-yellow-600'
                        }`}>
                          {tx.status}
                        </span>
                      </td>
                      <td className="py-4 text-gray-500 text-sm hidden md:table-cell">
                        {new Date(tx.timestamp).toLocaleDateString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="text-center py-8">
              <div className="w-16 h-16 bg-gray-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
                <span className="text-3xl">📭</span>
              </div>
              <h3 className="text-gray-800 font-semibold mb-2">No Transactions Yet</h3>
              <p className="text-gray-600">Start by sending or receiving some coins!</p>
            </div>
          )}
        </div>


      </div>
    </div>
  );
};

export default Dashboard;
