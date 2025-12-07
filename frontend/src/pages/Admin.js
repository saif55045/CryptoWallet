import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Navbar from '../components/Navbar';
import api from '../services/api';
import { useAuth } from '../context/AuthContext';

const Admin = () => {
  const navigate = useNavigate();
  const { user, token } = useAuth();
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState(null);
  const [users, setUsers] = useState([]);
  const [transactions, setTransactions] = useState([]);
  const [blocks, setBlocks] = useState([]);
  const [activeTab, setActiveTab] = useState('dashboard');
  const [error, setError] = useState('');

  useEffect(() => {
    if (!user?.isAdmin) {
      navigate('/dashboard');
      return;
    }
    fetchData();
  }, [user, navigate]);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [statsRes, usersRes, txRes, blocksRes] = await Promise.all([
        api.admin.getStats(token),
        api.admin.getUsers(token),
        api.admin.getTransactions(token),
        api.admin.getBlocks(token),
      ]);
      setStats(statsRes.data);
      setUsers(usersRes.data.users || []);
      setTransactions(txRes.data.transactions || []);
      setBlocks(blocksRes.data.blocks || []);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch admin data');
      if (err.response?.status === 403) {
        navigate('/dashboard');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleToggleAdmin = async (userId) => {
    try {
      await api.admin.toggleAdmin(token, userId);
      fetchData();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to toggle admin status');
    }
  };

  const handleDeleteUser = async (userId) => {
    if (!window.confirm('Are you sure you want to delete this user?')) return;
    try {
      await api.admin.deleteUser(token, userId);
      fetchData();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to delete user');
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString();
  };

  if (loading) {
    return (
      <div className="min-h-screen">
        <Navbar />
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <div className="w-16 h-16 border-4 border-indigo-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
            <p className="text-gray-600 mt-4">Loading admin panel...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-800">🛡️ Admin Panel</h1>
          <p className="text-gray-600 mt-2">System administration and monitoring</p>
        </div>

        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-xl">
            {error}
          </div>
        )}

        {/* Navigation Tabs */}
        <div className="flex gap-2 mb-6 bg-gray-100 p-1 rounded-xl w-fit">
          {['dashboard', 'users', 'transactions', 'blocks'].map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                activeTab === tab
                  ? 'bg-white text-indigo-600 shadow'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              {tab.charAt(0).toUpperCase() + tab.slice(1)}
            </button>
          ))}
        </div>

        {/* Dashboard Tab */}
        {activeTab === 'dashboard' && stats && (
          <div className="space-y-6">
            {/* Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <div className="glass-card rounded-2xl p-6">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
                    <span className="text-2xl">👥</span>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">Total Users</p>
                    <p className="text-2xl font-bold text-gray-800">{stats.users?.total || 0}</p>
                    <p className="text-xs text-green-600">{stats.users?.verified || 0} verified</p>
                  </div>
                </div>
              </div>

              <div className="glass-card rounded-2xl p-6">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center">
                    <span className="text-2xl">💼</span>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">Total Wallets</p>
                    <p className="text-2xl font-bold text-gray-800">{stats.wallets || 0}</p>
                  </div>
                </div>
              </div>

              <div className="glass-card rounded-2xl p-6">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center">
                    <span className="text-2xl">📊</span>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">Transactions</p>
                    <p className="text-2xl font-bold text-gray-800">{stats.transactions?.total || 0}</p>
                    <p className="text-xs text-yellow-600">{stats.transactions?.pending || 0} pending</p>
                  </div>
                </div>
              </div>

              <div className="glass-card rounded-2xl p-6">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 bg-amber-100 rounded-xl flex items-center justify-center">
                    <span className="text-2xl">⛏️</span>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">Blocks Mined</p>
                    <p className="text-2xl font-bold text-gray-800">{stats.blocks || 0}</p>
                  </div>
                </div>
              </div>
            </div>

            {/* More Stats */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="glass-card rounded-2xl p-6">
                <h3 className="text-lg font-semibold text-gray-800 mb-4">💰 System Balance</h3>
                <p className="text-3xl font-bold gradient-text">{stats.totalBalance?.toFixed(2) || '0.00'}</p>
                <p className="text-sm text-gray-500">Total coins in circulation</p>
              </div>

              <div className="glass-card rounded-2xl p-6">
                <h3 className="text-lg font-semibold text-gray-800 mb-4">📦 UTXOs</h3>
                <p className="text-3xl font-bold text-green-600">{stats.utxos?.unspent || 0}</p>
                <p className="text-sm text-gray-500">Unspent outputs ({stats.utxos?.total || 0} total)</p>
              </div>

              <div className="glass-card rounded-2xl p-6">
                <h3 className="text-lg font-semibold text-gray-800 mb-4">🕌 Zakat Payments</h3>
                <p className="text-3xl font-bold text-purple-600">{stats.zakatPayments || 0}</p>
                <p className="text-sm text-gray-500">Total zakat transactions</p>
              </div>
            </div>
          </div>
        )}

        {/* Users Tab */}
        {activeTab === 'users' && (
          <div className="glass-card rounded-2xl p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">All Users ({users.length})</h3>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-200">
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">User</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Email</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Wallet ID</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Status</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Joined</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((u) => (
                    <tr key={u.id} className="border-b border-gray-100 hover:bg-gray-50">
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 bg-indigo-100 rounded-full flex items-center justify-center">
                            <span className="text-indigo-600 font-semibold">{u.fullName?.charAt(0) || '?'}</span>
                          </div>
                          <div>
                            <p className="font-medium text-gray-800">{u.fullName}</p>
                            {u.isAdmin && <span className="text-xs text-indigo-600 font-medium">Admin</span>}
                          </div>
                        </div>
                      </td>
                      <td className="py-3 px-4 text-sm text-gray-600">{u.email}</td>
                      <td className="py-3 px-4 text-sm text-gray-500 font-mono truncate max-w-[150px]">
                        {u.walletId || 'No wallet'}
                      </td>
                      <td className="py-3 px-4">
                        <span className={`px-2 py-1 text-xs font-medium rounded-full ${
                          u.isVerified ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'
                        }`}>
                          {u.isVerified ? 'Verified' : 'Pending'}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-sm text-gray-500">{formatDate(u.createdAt)}</td>
                      <td className="py-3 px-4">
                        <div className="flex gap-2">
                          <button
                            onClick={() => handleToggleAdmin(u.id)}
                            className={`px-3 py-1 text-xs font-medium rounded-lg transition-colors ${
                              u.isAdmin
                                ? 'bg-yellow-100 text-yellow-700 hover:bg-yellow-200'
                                : 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200'
                            }`}
                          >
                            {u.isAdmin ? 'Remove Admin' : 'Make Admin'}
                          </button>
                          <button
                            onClick={() => handleDeleteUser(u.id)}
                            className="px-3 py-1 text-xs font-medium rounded-lg bg-red-100 text-red-700 hover:bg-red-200 transition-colors"
                          >
                            Delete
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {users.length === 0 && (
                <p className="text-center text-gray-500 py-8">No users found</p>
              )}
            </div>
          </div>
        )}

        {/* Transactions Tab */}
        {activeTab === 'transactions' && (
          <div className="glass-card rounded-2xl p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Recent Transactions ({transactions.length})</h3>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-200">
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">TX ID</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Type</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Amount</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">From</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Status</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Time</th>
                  </tr>
                </thead>
                <tbody>
                  {transactions.map((tx, index) => (
                    <tr key={index} className="border-b border-gray-100 hover:bg-gray-50">
                      <td className="py-3 px-4 text-sm text-gray-500 font-mono truncate max-w-[100px]">
                        {tx.transactionId?.substring(0, 12)}...
                      </td>
                      <td className="py-3 px-4">
                        <span className={`px-2 py-1 text-xs font-medium rounded-full ${
                          tx.type === 'coinbase' ? 'bg-amber-100 text-amber-700' :
                          tx.type === 'zakat' ? 'bg-purple-100 text-purple-700' :
                          'bg-blue-100 text-blue-700'
                        }`}>
                          {tx.type}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-sm font-medium text-gray-900">{tx.totalOutput?.toFixed(2) || '0.00'}</td>
                      <td className="py-3 px-4 text-sm text-gray-500 font-mono truncate max-w-[100px]">
                        {tx.senderWallet?.substring(0, 12) || 'System'}...
                      </td>
                      <td className="py-3 px-4">
                        <span className={`px-2 py-1 text-xs font-medium rounded-full ${
                          tx.status === 'confirmed' ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'
                        }`}>
                          {tx.status}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-sm text-gray-500">{formatDate(tx.timestamp)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {transactions.length === 0 && (
                <p className="text-center text-gray-500 py-8">No transactions found</p>
              )}
            </div>
          </div>
        )}

        {/* Blocks Tab */}
        {activeTab === 'blocks' && (
          <div className="glass-card rounded-2xl p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Blockchain ({blocks.length} blocks)</h3>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-200">
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Index</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Hash</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Merkle Root</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Miner</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Reward</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">TXs</th>
                    <th className="text-left py-3 px-4 text-sm font-semibold text-gray-600">Time</th>
                  </tr>
                </thead>
                <tbody>
                  {blocks.map((block, index) => (
                    <tr key={index} className="border-b border-gray-100 hover:bg-gray-50">
                      <td className="py-3 px-4 text-sm font-medium text-gray-900">#{block.index}</td>
                      <td className="py-3 px-4 text-sm text-gray-500 font-mono truncate max-w-[100px]">
                        {block.hash?.substring(0, 12)}...
                      </td>
                      <td className="py-3 px-4 text-sm text-gray-500 font-mono truncate max-w-[100px]">
                        {block.merkleRoot?.substring(0, 12) || 'N/A'}...
                      </td>
                      <td className="py-3 px-4 text-sm text-gray-500 font-mono truncate max-w-[100px]">
                        {block.minerWalletId?.substring(0, 12) || 'Genesis'}...
                      </td>
                      <td className="py-3 px-4 text-sm font-medium text-amber-600">{block.miningReward?.toFixed(2) || '0.00'}</td>
                      <td className="py-3 px-4 text-sm text-gray-600">{block.transactionCount || 0}</td>
                      <td className="py-3 px-4 text-sm text-gray-500">{formatDate(block.timestamp)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {blocks.length === 0 && (
                <p className="text-center text-gray-500 py-8">No blocks found</p>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Admin;
