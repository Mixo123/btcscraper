import { useEffect, useState } from 'react';
import { View, Text, ActivityIndicator } from 'react-native'
import { SafeAreaView } from 'react-native-safe-area-context';
import axios from 'axios'

interface BTCProps {
  title: string;
  price: number;
}

const HomeScreen: React.FC = () => {
  const [data, setData] = useState<BTCProps[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<boolean>(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('http://localhost:8080/btc');
        setData(response.data);
      } catch (error) {
        console.error(error);
        setError(true);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return <ActivityIndicator size="large" color="#0000ff" />;
  }

  if (error) {
    return <Text>Error occurred while fetching data.</Text>;
  }
  
  return (
    <SafeAreaView>
      <View className='flex-1 justify-center items-center bg-white'>
        {data.map((item, index) => (
          <View key={index} className='flex flex-col gap-5'>
            <Text className='text-black text-3xl font-bold'>{item.title}</Text>
            <Text className='text-black text-3xl font-bold'>{item.price}</Text>
          </View>
        ))}
      </View>
    </SafeAreaView>
  )
}

export default HomeScreen