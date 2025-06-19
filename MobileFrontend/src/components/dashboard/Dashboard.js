import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  ScrollView,
  RefreshControl,
  StyleSheet,
  SafeAreaView,
  TouchableOpacity,
} from 'react-native';
import PortfolioPerformance from './PortfolioPerformance';
import Positions from './Positions';
import RecentOrders from './RecentOrders';
import MarketData from './MarketData';

const Dashboard = ({ onOrderComplete }) => {
  const [refreshing, setRefreshing] = useState(false);
  const [key, setKey] = useState(0);

  const onRefresh = () => {
    setRefreshing(true);
    setKey(prev => prev + 1);
    setTimeout(() => {
      setRefreshing(false);
    }, 1000);
  };

  useEffect(() => {
    if (onOrderComplete) {
      setKey(prev => prev + 1);
    }
  }, [onOrderComplete]);

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView
        style={styles.scrollView}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
        showsVerticalScrollIndicator={false}
      >
        <View style={styles.content}>
          {/* Portfolio Performance Section */}
          <View style={styles.section}>
            <PortfolioPerformance key={`portfolio-${key}`} />
          </View>

          {/* Quick Stats Row */}
          <View style={styles.statsRow}>
            <View style={styles.statCard}>
              <Positions key={`positions-${key}`} isCompact={true} />
            </View>
          </View>

          {/* Recent Orders Section */}
          <View style={styles.section}>
            <RecentOrders 
              key={`orders-${key}`} 
              onOrderComplete={onOrderComplete}
              isMobile={true}
            />
          </View>

          {/* Market Data Section */}
          <View style={styles.section}>
            <MarketData key={`market-${key}`} />
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f9fafb',
  },
  scrollView: {
    flex: 1,
  },
  content: {
    padding: 16,
    paddingBottom: 32,
  },
  section: {
    marginBottom: 20,
    backgroundColor: 'white',
    borderRadius: 16,
    padding: 16,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.1,
    shadowRadius: 3.84,
    elevation: 5,
  },
  statsRow: {
    flexDirection: 'row',
    marginBottom: 20,
  },
  statCard: {
    flex: 1,
    backgroundColor: 'white',
    borderRadius: 16,
    padding: 16,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.1,
    shadowRadius: 3.84,
    elevation: 5,
  },
});

export default Dashboard; 