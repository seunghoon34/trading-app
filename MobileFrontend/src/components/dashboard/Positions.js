import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ActivityIndicator,
  FlatList,
  TouchableOpacity,
} from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';

const Positions = ({ isCompact = false }) => {
  const [positions, setPositions] = useState([]);
  const [loading, setLoading] = useState(true);

  const API_BASE_URL = __DEV__ 
    ? 'http://localhost:3000' 
    : 'https://your-production-api.com';

  useEffect(() => {
    fetchPositions();
  }, []);

  const fetchPositions = async () => {
    try {
      const token = await AsyncStorage.getItem('authToken');
      const response = await fetch(`${API_BASE_URL}/api/v1/portfolio/positions`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        const data = await response.json();
        setPositions(data.positions || []);
      }
    } catch (error) {
      console.error('Error fetching positions:', error);
    } finally {
      setLoading(false);
    }
  };

  // Mock data for demonstration
  const mockPositions = [
    {
      symbol: 'AAPL',
      name: 'Apple Inc.',
      quantity: 50,
      currentPrice: 185.25,
      marketValue: 9262.50,
      unrealizedPL: 1125.00,
      unrealizedPLPercent: 13.82,
      averageCost: 162.75,
    },
    {
      symbol: 'TSLA',
      name: 'Tesla Inc.',
      quantity: 25,
      currentPrice: 248.50,
      marketValue: 6212.50,
      unrealizedPL: -875.00,
      unrealizedPLPercent: -12.33,
      averageCost: 283.50,
    },
    {
      symbol: 'NVDA',
      name: 'NVIDIA Corporation',
      quantity: 30,
      currentPrice: 876.50,
      marketValue: 26295.00,
      unrealizedPL: 3450.00,
      unrealizedPLPercent: 15.12,
      averageCost: 761.50,
    },
  ];

  const displayPositions = positions.length > 0 ? positions : mockPositions;

  const renderPosition = ({ item }) => (
    <TouchableOpacity style={styles.positionItem}>
      <View style={styles.positionHeader}>
        <View style={styles.symbolContainer}>
          <Text style={styles.symbol}>{item.symbol}</Text>
          {!isCompact && (
            <Text style={styles.companyName} numberOfLines={1}>
              {item.name}
            </Text>
          )}
        </View>
        <View style={styles.valueContainer}>
          <Text style={styles.marketValue}>
            ${item.marketValue.toLocaleString('en-US', { minimumFractionDigits: 2 })}
          </Text>
          <Text style={[
            styles.plValue,
            item.unrealizedPL >= 0 ? styles.positive : styles.negative
          ]}>
            {item.unrealizedPL >= 0 ? '+' : ''}${Math.abs(item.unrealizedPL).toLocaleString('en-US', { minimumFractionDigits: 2 })}
          </Text>
        </View>
      </View>
      
      {!isCompact && (
        <View style={styles.positionDetails}>
          <Text style={styles.detailText}>
            {item.quantity} shares @ ${item.currentPrice.toFixed(2)}
          </Text>
          <Text style={[
            styles.percentChange,
            item.unrealizedPL >= 0 ? styles.positive : styles.negative
          ]}>
            {item.unrealizedPL >= 0 ? '+' : ''}{item.unrealizedPLPercent.toFixed(2)}%
          </Text>
        </View>
      )}
    </TouchableOpacity>
  );

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="small" color="#7c3aed" />
        <Text style={styles.loadingText}>Loading positions...</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>
          {isCompact ? 'Top Positions' : 'Your Positions'}
        </Text>
        {displayPositions.length > 0 && (
          <Text style={styles.totalPositions}>
            {displayPositions.length} position{displayPositions.length !== 1 ? 's' : ''}
          </Text>
        )}
      </View>

      {displayPositions.length === 0 ? (
        <View style={styles.emptyContainer}>
          <Text style={styles.emptyText}>No positions found</Text>
          <Text style={styles.emptySubtext}>Start investing to see your positions here</Text>
        </View>
      ) : (
        <FlatList
          data={isCompact ? displayPositions.slice(0, 3) : displayPositions}
          renderItem={renderPosition}
          keyExtractor={(item) => item.symbol}
          scrollEnabled={!isCompact}
          showsVerticalScrollIndicator={false}
        />
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: 100,
  },
  loadingText: {
    marginTop: 8,
    fontSize: 14,
    color: '#6b7280',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  title: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#1f2937',
  },
  totalPositions: {
    fontSize: 14,
    color: '#6b7280',
  },
  positionItem: {
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#f3f4f6',
  },
  positionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
  },
  symbolContainer: {
    flex: 1,
  },
  symbol: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#1f2937',
  },
  companyName: {
    fontSize: 12,
    color: '#6b7280',
    marginTop: 2,
  },
  valueContainer: {
    alignItems: 'flex-end',
  },
  marketValue: {
    fontSize: 16,
    fontWeight: '600',
    color: '#1f2937',
  },
  plValue: {
    fontSize: 14,
    fontWeight: '600',
    marginTop: 2,
  },
  positionDetails: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 8,
  },
  detailText: {
    fontSize: 12,
    color: '#6b7280',
  },
  percentChange: {
    fontSize: 12,
    fontWeight: '600',
  },
  positive: {
    color: '#10b981',
  },
  negative: {
    color: '#ef4444',
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: 40,
  },
  emptyText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#6b7280',
    marginBottom: 4,
  },
  emptySubtext: {
    fontSize: 14,
    color: '#9ca3af',
    textAlign: 'center',
  },
});

export default Positions; 