import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ActivityIndicator,
  FlatList,
  TouchableOpacity,
  Alert,
} from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';

const RecentOrders = ({ onOrderComplete, isMobile = false }) => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);

  const API_BASE_URL = __DEV__ 
    ? 'http://localhost:3000' 
    : 'https://your-production-api.com';

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    try {
      const token = await AsyncStorage.getItem('authToken');
      const response = await fetch(`${API_BASE_URL}/api/v1/trading/orders`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        const data = await response.json();
        setOrders(data.orders || []);
      }
    } catch (error) {
      console.error('Error fetching orders:', error);
    } finally {
      setLoading(false);
    }
  };

  // Mock data for demonstration
  const mockOrders = [
    {
      id: '1',
      symbol: 'AAPL',
      side: 'buy',
      quantity: 10,
      price: 185.25,
      status: 'filled',
      createdAt: '2024-06-19T10:30:00Z',
    },
    {
      id: '2',
      symbol: 'TSLA',
      side: 'sell',
      quantity: 5,
      price: 248.50,
      status: 'pending',
      createdAt: '2024-06-19T09:15:00Z',
    },
  ];

  const displayOrders = orders.length > 0 ? orders : mockOrders;

  const renderOrder = ({ item }) => (
    <View style={styles.orderItem}>
      <View style={styles.orderHeader}>
        <Text style={styles.symbol}>{item.symbol}</Text>
        <Text style={[styles.side, item.side === 'buy' ? styles.buy : styles.sell]}>
          {item.side.toUpperCase()}
        </Text>
      </View>
      <Text style={styles.details}>
        {item.quantity} shares @ ${item.price.toFixed(2)}
      </Text>
      <Text style={styles.status}>Status: {item.status}</Text>
    </View>
  );

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="small" color="#7c3aed" />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Recent Orders</Text>
      <FlatList
        data={displayOrders.slice(0, 5)}
        renderItem={renderOrder}
        keyExtractor={(item) => item.id}
        showsVerticalScrollIndicator={false}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  loadingContainer: {
    justifyContent: 'center',
    alignItems: 'center',
    height: 100,
  },
  title: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  orderItem: {
    padding: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  orderHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  symbol: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  side: {
    fontSize: 14,
    fontWeight: '600',
  },
  buy: {
    color: '#10b981',
  },
  sell: {
    color: '#ef4444',
  },
  details: {
    fontSize: 14,
    color: '#666',
    marginTop: 4,
  },
  status: {
    fontSize: 12,
    color: '#666',
    marginTop: 4,
  },
});

export default RecentOrders; 