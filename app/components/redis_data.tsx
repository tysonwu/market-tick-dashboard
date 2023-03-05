import { useState, useEffect } from 'react';
import Redis, {RedisOptions} from 'ioredis';
import type { GetServerSideProps } from 'next';

interface RedisClient {
  subscribe: (channel: string) => void;
  unsubscribe: (channel: string) => void;
  on: (event: string, listener: (channel: string, message: string) => void) => void;
}

interface TickData {
  Symbol: string;
  Price: number;
  Exchange: string;
  Time: Date;
}

type Props = {
  data: TickData | null;
}

// TODO: redisClient cannot live outside SSR
// but we still want to subscribe the change in key value and reflect it on this UI
export default function RedisData({ data }: Props) {
  const [tickData, setData] = useState<TickData | null>(data);
  const k: string = 'ticks';

  useEffect(() => {
    const dbConfig: RedisOptions = {
      host: 'localhost',
      port: 6379,
      db: 0,
    };
    const redisClient: RedisClient = new Redis(dbConfig);
    redisClient.subscribe(k); // the key name in redis
    redisClient.on('message', (_, message) => {
      const newTick = JSON.parse(message) as TickData;
      setData(newTick);
    });

    return () => {
      redisClient.unsubscribe(k);
    };
  }, []);

  return (
    <div>
      {data ? (
        <ul>
          <li>Symbol: {data.Symbol}</li>
          <li>Price: {data.Price}</li>
          <li>Exchange: {data.Exchange}</li>
          <li>Time: {data.Time.toString()}</li>
        </ul>
      ) : (
        <p>No data available</p>
      )}
    </div>
  );
}

export const getServerSideProps: GetServerSideProps<Props> = async () => {
  const dbConfig: RedisOptions = {
    host: 'localhost',
    port: 6379,
    db: 0,
  };
  const redisClient: RedisClient = new Redis(dbConfig);
  redisClient.subscribe(k); // the key name in redis
  redisClient.on('message', (_, message) => {
    const newTick = JSON.parse(message) as TickData;
    setData(newTick);
  });

  return {
    props: {
      data: myData ?? null,
    },
  };
};