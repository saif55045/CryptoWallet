import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
import api from '../services/api';
import Navbar from '../components/Navbar';

const Zakat = () => {
  const { token } = useAuth();
  const [summary, setSummary] = useState(null);
  const [settings, setSettings] = useState(null);
  const [history, setHistory] = useState({ calculations: [], payments: [] });
  const [loading, setLoading] = useState(true);
  const [calculating, setCalculating] = useState(false);
  const [paying, setPaying] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [activeTab, setActiveTab] = useState('overview');

  const fetchData = useCallback(async () => {
    try {
      const [settingsRes, summaryRes, historyRes] = await Promise.all([
        api.zakat.getSettings(),
        api.zakat.getSummary(token),
        api.zakat.getHistory(token)
      ]);
      
      setSettings(settingsRes.data.settings);
      setSummary(summaryRes.data.summary);
      setHistory(historyRes.data);
    } catch (err) {
      if (err.response?.status !== 404) {
        setError('Failed to load zakat data');
      }
    } finally {
      setLoading(false);
    }
  }, [token]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleCalculate = async () => {
    setCalculating(true);
    setError('');
    setSuccess('');
    
    try {
      await api.zakat.calculate(token);
      setSuccess('Zakat calculated successfully!');
      fetchData();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to calculate zakat');
    } finally {
      setCalculating(false);
    }
  };

  const handlePay = async () => {
    setPaying(true);
    setError('');
    setSuccess('');
    
    try {
      const amount = summary?.zakatDue || 0;
      if (amount <= 0) {
        setError('No zakat due to pay');
        return;
      }
      
      await api.zakat.pay(token, {
        amount,
        calculationId: summary?.lastCalculation?._id
      });
      
      setSuccess('Zakat payment submitted! It will be confirmed when mined.');
      fetchData();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to pay zakat');
    } finally {
      setPaying(false);
    }
  };

  const isEligible = summary?.currentBalance >= (settings?.nisabThreshold || 1000);
  const zakatDue = summary?.zakatDue || 0;

  if (loading) {
    return (
      <div className="min-h-screen">
        <Navbar />
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <div className="w-16 h-16 border-4 border-yellow-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
            <p className="text-gray-600 mt-4">Loading zakat information...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <Navbar />
      
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold gradient-text mb-2">🕌 Zakat Management</h1>
          <p className="text-gray-600">Calculate and pay your zakat obligations</p>
        </div>

        {error && (
          <div className="mb-6 bg-red-500/10 border border-red-500/50 text-red-400 px-4 py-3 rounded-xl flex items-center gap-2">
            <svg className="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            {error}
          </div>
        )}

        {success && (
          <div className="mb-6 bg-green-500/10 border border-green-500/50 text-green-400 px-4 py-3 rounded-xl flex items-center gap-2">
            <svg className="w-5 h-5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
            </svg>
            {success}
          </div>
        )}

        {/* Tabs */}
        <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
          {['overview', 'calculate', 'history', 'info'].map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`px-4 py-2 rounded-xl font-medium transition-all capitalize whitespace-nowrap ${
                activeTab === tab
                  ? 'bg-yellow-500 text-black'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {tab}
            </button>
          ))}
        </div>

        {activeTab === 'overview' && (
          <div className="space-y-6">
            {/* Status Banner */}
            <div className={`glass-card rounded-2xl p-8 ${
              isEligible ? 'bg-gradient-to-r from-yellow-500/10 to-orange-500/10 border-yellow-500/20' : 'bg-gradient-to-r from-green-500/10 to-emerald-500/10 border-green-500/20'
            }`}>
              <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-6">
                <div>
                  <div className="flex items-center gap-3 mb-2">
                    <span className="text-4xl">{isEligible ? '🕌' : '✅'}</span>
                    <div>
                      <h2 className="text-2xl font-bold text-gray-800">
                        {isEligible ? 'Zakat is Due' : 'No Zakat Due'}
                      </h2>
                      <p className="text-gray-600">
                        {isEligible 
                          ? `Your balance exceeds the nisab threshold of ${settings?.nisabThreshold || 1000} coins`
                          : 'Your balance is below the nisab threshold'}
                      </p>
                    </div>
                  </div>
                </div>
                
                {isEligible && zakatDue > 0 && (
                  <div className="text-right">
                    <p className="text-gray-600 text-sm">Amount Due</p>
                    <p className="text-4xl font-bold text-yellow-600">{zakatDue.toFixed(2)}</p>
                    <p className="text-gray-600 text-sm">coins ({settings?.zakatRate || 2.5}%)</p>
                  </div>
                )}
              </div>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
              <div className="glass-card rounded-2xl p-6">
                <p className="text-gray-600 text-sm mb-1">Current Balance</p>
                <p className="text-2xl font-bold text-gray-800">{summary?.currentBalance?.toFixed(2) || '0.00'}</p>
              </div>
              <div className="glass-card rounded-2xl p-6">
                <p className="text-gray-600 text-sm mb-1">Nisab Threshold</p>
                <p className="text-2xl font-bold text-indigo-600">{settings?.nisabThreshold || 1000}</p>
              </div>
              <div className="glass-card rounded-2xl p-6">
                <p className="text-gray-600 text-sm mb-1">Zakat Rate</p>
                <p className="text-2xl font-bold text-purple-600">{settings?.zakatRate || 2.5}%</p>
              </div>
              <div className="glass-card rounded-2xl p-6">
                <p className="text-gray-600 text-sm mb-1">Total Paid</p>
                <p className="text-2xl font-bold text-green-600">{summary?.totalPaid?.toFixed(2) || '0.00'}</p>
              </div>
            </div>

            {/* Quick Actions */}
            {isEligible && (
              <div className="glass-card rounded-2xl p-6">
                <h3 className="text-lg font-semibold text-gray-800 mb-4">Quick Actions</h3>
                <div className="flex flex-wrap gap-4">
                  <button
                    onClick={handleCalculate}
                    disabled={calculating}
                    className="px-6 py-3 rounded-xl font-semibold text-white transition-all duration-300 transform hover:scale-105 flex items-center gap-2 disabled:opacity-50"
                    style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
                  >
                    {calculating ? (
                      <>
                        <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                        </svg>
                        Calculating...
                      </>
                    ) : (
                      <>
                        <span>🧮</span>
                        Calculate Zakat
                      </>
                    )}
                  </button>

                  {zakatDue > 0 && (
                    <button
                      onClick={handlePay}
                      disabled={paying}
                      className="px-6 py-3 rounded-xl font-semibold text-black transition-all duration-300 transform hover:scale-105 flex items-center gap-2 disabled:opacity-50"
                      style={{ background: 'linear-gradient(135deg, #f2c94c 0%, #f2994a 100%)' }}
                    >
                      {paying ? (
                        <>
                          <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                          </svg>
                          Processing...
                        </>
                      ) : (
                        <>
                          <span>💝</span>
                          Pay Zakat ({zakatDue.toFixed(2)})
                        </>
                      )}
                    </button>
                  )}
                </div>
              </div>
            )}
          </div>
        )}

        {activeTab === 'calculate' && (
          <div className="glass-card rounded-2xl p-6">
            <h2 className="text-xl font-bold text-gray-800 mb-6">🧮 Zakat Calculator</h2>
            
            <div className="space-y-6">
              <div className="bg-gray-50 rounded-xl p-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <p className="text-gray-600 text-sm mb-1">Your Balance</p>
                    <p className="text-3xl font-bold text-gray-800">{summary?.currentBalance?.toFixed(4) || '0.00'}</p>
                  </div>
                  <div>
                    <p className="text-gray-600 text-sm mb-1">Nisab (Minimum)</p>
                    <p className="text-3xl font-bold text-indigo-600">{settings?.nisabThreshold || 1000}</p>
                  </div>
                </div>
              </div>

              <div className="flex items-center gap-4 p-4 rounded-xl bg-gray-50">
                <span className="text-3xl">{isEligible ? '✅' : '❌'}</span>
                <div>
                  <p className="text-gray-800 font-semibold">
                    {isEligible ? 'Zakat is Obligatory' : 'Zakat Not Required'}
                  </p>
                  <p className="text-gray-600 text-sm">
                    {isEligible 
                      ? 'Your wealth exceeds the nisab threshold'
                      : 'Build your wealth above the nisab to be eligible'}
                  </p>
                </div>
              </div>

              {isEligible && (
                <div className="bg-gradient-to-r from-yellow-500/10 to-orange-500/10 rounded-xl p-6 border border-yellow-500/20">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-gray-600 text-sm">Zakat Amount ({settings?.zakatRate || 2.5}%)</p>
                      <p className="text-4xl font-bold text-yellow-600">{zakatDue.toFixed(4)}</p>
                    </div>
                    <button
                      onClick={handleCalculate}
                      disabled={calculating}
                      className="btn-primary flex items-center gap-2"
                    >
                      {calculating ? 'Calculating...' : 'Recalculate'}
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}

        {activeTab === 'history' && (
          <div className="space-y-6">
            {/* Payment History */}
            <div className="glass-card rounded-2xl p-6">
              <h2 className="text-xl font-bold text-gray-800 mb-6">💝 Payment History</h2>
              
              {history.payments?.length > 0 ? (
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="text-left text-gray-600 text-sm border-b border-gray-200">
                        <th className="pb-4 font-medium">Date</th>
                        <th className="pb-4 font-medium">Amount</th>
                        <th className="pb-4 font-medium">Status</th>
                        <th className="pb-4 font-medium">Transaction</th>
                      </tr>
                    </thead>
                    <tbody>
                      {history.payments.map((payment, index) => (
                        <tr key={index} className="border-b border-gray-100">
                          <td className="py-4 text-gray-700">
                            {new Date(payment.paidAt).toLocaleDateString()}
                          </td>
                          <td className="py-4 text-green-400 font-semibold">
                            {payment.amount?.toFixed(2)}
                          </td>
                          <td className="py-4">
                            <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                              payment.status === 'confirmed'
                                ? 'bg-green-500/20 text-green-400'
                                : 'bg-yellow-500/20 text-yellow-400'
                            }`}>
                              {payment.status}
                            </span>
                          </td>
                          <td className="py-4">
                            <code className="text-gray-500 text-xs">{payment.transactionId?.substring(0, 16)}...</code>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div className="text-center py-8 text-gray-600">
                  <span className="text-4xl mb-4 block">📭</span>
                  <p>No payment history yet</p>
                </div>
              )}
            </div>

            {/* Calculation History */}
            <div className="glass-card rounded-2xl p-6">
              <h2 className="text-xl font-bold text-gray-800 mb-6">🧮 Calculation History</h2>
              
              {history.calculations?.length > 0 ? (
                <div className="space-y-3">
                  {history.calculations.map((calc, index) => (
                    <div key={index} className="bg-gray-50 rounded-xl p-4 flex items-center justify-between">
                      <div>
                        <p className="text-gray-800 font-medium">
                          Balance: {calc.eligibleBalance?.toFixed(2)} → Zakat: {calc.zakatAmount?.toFixed(2)}
                        </p>
                        <p className="text-gray-600 text-sm">
                          {new Date(calc.calculatedAt).toLocaleString()}
                        </p>
                      </div>
                      <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                        calc.isPaid
                          ? 'bg-green-500/20 text-green-400'
                          : 'bg-yellow-500/20 text-yellow-400'
                      }`}>
                        {calc.isPaid ? 'Paid' : 'Pending'}
                      </span>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8 text-gray-600">
                  <span className="text-4xl mb-4 block">📊</span>
                  <p>No calculations yet</p>
                </div>
              )}
            </div>
          </div>
        )}

        {activeTab === 'info' && (
          <div className="space-y-6">
            <div className="glass-card rounded-2xl p-6">
              <h2 className="text-xl font-bold text-gray-800 mb-6">📖 About Zakat</h2>
              
              <div className="prose max-w-none">
                <p className="text-gray-700 mb-4">
                  Zakat is one of the Five Pillars of Islam. It is a form of obligatory charity that is 
                  required of Muslims who meet certain wealth thresholds.
                </p>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-6">
                  <div className="bg-gray-50 rounded-xl p-4">
                    <h4 className="text-gray-800 font-semibold mb-2">📊 Rate</h4>
                    <p className="text-gray-600">2.5% of eligible wealth annually</p>
                  </div>
                  <div className="bg-gray-50 rounded-xl p-4">
                    <h4 className="text-gray-800 font-semibold mb-2">💰 Nisab</h4>
                    <p className="text-gray-600">Minimum wealth threshold to qualify</p>
                  </div>
                </div>
              </div>
            </div>

            <div className="glass-card rounded-2xl p-6">
              <h2 className="text-xl font-bold text-gray-800 mb-6">👥 Eligible Recipients</h2>
              
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                {[
                  { name: 'The Poor', icon: '🏚️' },
                  { name: 'The Needy', icon: '🤲' },
                  { name: 'Zakat Collectors', icon: '📋' },
                  { name: 'New Muslims', icon: '🕌' },
                  { name: 'Freeing Captives', icon: '🔓' },
                  { name: 'Debtors', icon: '📝' },
                  { name: "In Allah's Cause", icon: '⚔️' },
                  { name: 'Travelers', icon: '🧳' },
                ].map((recipient, index) => (
                  <div key={index} className="bg-gray-50 rounded-xl p-4 text-center">
                    <span className="text-2xl block mb-2">{recipient.icon}</span>
                    <p className="text-gray-700 text-sm">{recipient.name}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Zakat;
