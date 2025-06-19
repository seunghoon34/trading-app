import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ActivityIndicator,
  FlatList,
} from 'react-native';

const MarketData = () => {
  const [marketData, setMarketData] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Mock data for demonstration
    setTimeout(() => {
      const mockData = [
        {
          symbol: 'AAPL',
          name: 'Apple Inc.',
          price: 185.25,
          change: 2.35,
          changePercent: 1.29,
        },
        {
          symbol: 'TSLA',
          name: 'Tesla Inc.',
          price: 248.50,
          change: -12.25,
          changePercent: -4.70,
        },
        {
          symbol: 'NVDA',
          name: 'NVIDIA Corporation',
          price: 876.50,
          change: 34.75,
          changePercent: 4.13,
        },
      ];
      
      setMarketData(mockData);
      setLoading(false);
    }, 500);
  }, []);

  const renderStock = ({ item }) => (
    <View style={styles.stockItem}>
      <View style={styles.stockInfo}>
        <Text style={styles.symbol}>{item.symbol}</Text>
        <Text style={styles.price}>${item.price.toFixed(2)}</Text>
      </View>
      <View style={styles.changeInfo}>
        <Text style={[
          styles.change,
          item.change >= 0 ? styles.positive : styles.negative
        ]}>
          {item.change >= 0 ? '+' : ''}${Math.abs(item.change).toFixed(2)}
        </Text>
        <Text style={[
          styles.changePercent,
          item.change >= 0 ? styles.positive : styles.negative
        ]}>
          ({item.change >= 0 ? '+' : ''}{item.changePercent.toFixed(2)}%)
        </Text>
      </View>
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
      <Text style={styles.title}>Market Data</Text>
      <FlatList
        data={marketData}
        renderItem={renderStock}
        keyExtractor={(item) => item.symbol}
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
  stockItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f0f0f0',
  },
  stockInfo: {
    flex: 1,
  },
  symbol: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  price: {
    fontSize: 14,
    color: '#666',
    marginTop: 2,
  },
  changeInfo: {
    alignItems: 'flex-end',
  },
  change: {
    fontSize: 14,
    fontWeight: '600',
  },
  changePercent: {
    fontSize: 12,
    marginTop: 2,
  },
  positive: {
    color: '#10b981',
  },
  negative: {
    color: '#ef4444',
  },
});

export default MarketData; 