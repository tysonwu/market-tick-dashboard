import { useState, useEffect, useRef } from 'react';
import { Table } from 'ka-table';
import { DataType, EditingMode, SortingMode } from 'ka-table/enums';

interface SymbolData {
  // StandardSymbol: string;
  symbol: string | null;
  // price: number | null;
  // AskPrice: number | null;
  // BidPrice: number | null;
  // Exchange: string;
  // Time: Date;
  binance: number | null;
  binanceBid: number | null;
  binanceAsk: number | null;
  kucoin: number | null;
  kucoinBid: number | null;
  kucoinAsk: number | null;
  id: number | null;
}

interface PriceTable { [key: string]: SymbolData }
// interface PriceTable extends Array<SymbolData>{}

const symbols: Array<string> = [
  'BTC-USDT',
  'ETH-USDT',
  'BNB-USDT',
  'XRP-USDT',
  'ADA-USDT',
  'DOGE-USDT',
  'MATIC-USDT',
  'SOL-USDT',
  'ETH-BTC',
]

const symbolDataDefaultValues: SymbolData = {
  symbol: null,
  // price: null,
  // AskPrice: null,
  // BidPrice: null,
  binance: null,
  binanceBid: null,
  binanceAsk: null,
  kucoin: null,
  kucoinBid: null,
  kucoinAsk: null,
  id: null,
};

const defaultPriceTable: PriceTable = {};
for (var i = 0; i < symbols.length; i++) {
  defaultPriceTable[symbols[i]] = structuredClone(symbolDataDefaultValues); // deepcopy!
  defaultPriceTable[symbols[i]].symbol = symbols[i];
  defaultPriceTable[symbols[i]].id = i;
}

export default function TickTable() {
  const tableRef = useRef(defaultPriceTable)
  const [table, setTable] = useState<PriceTable>(defaultPriceTable);

  useEffect(() => {
    tableRef.current = table
  }, [table]);

  useEffect(
    () => {
      var ws = new WebSocket('ws://localhost:8999');
      ws.onopen = () => {
        console.log('Connected to 8999');
      };
      ws.onmessage = e => {
        var { channel, message } = JSON.parse(e.data);
        var [ type, symbol, exchange ] = channel.split(':');
        var msg = JSON.parse(message);
        const newTable = structuredClone(tableRef.current);
        if (type === 'ticks') {
          if (exchange == 'binance') { newTable[symbol].binance = msg['Price']; };
          if (exchange == 'kucoin') { newTable[symbol].kucoin = msg['Price']; };
          setTable(newTable);
        };
        if (type === 'bidAskTicks') {
          if (exchange === 'binance') { newTable[symbol].binanceBid = msg['Bid']; newTable[symbol].binanceAsk = msg['Ask']; }
          if (exchange === 'kucoin') { newTable[symbol].kucoinBid = msg['Bid']; newTable[symbol].kucoinAsk = msg['Ask']; }
          setTable(newTable);
        };
      };
      return () => {
        ws.close();
      }
    },
    []
  );

  // const dataArray = Array(10).fill(0.12).map(
  //   (val, index) => ({
  //     symbol: `Symbol${index}`,
  //     binance: `B${index}`,
  //     kucoin: `K${index}`,
  //   }),
  // );
  // console.log(dataArray);

  return (
    <div>
      <Table
        columns={[
          { key: 'symbol', title: 'Symbol', dataType: DataType.String },
          { key: 'binanceAsk', title: 'Binance Bid', dataType: DataType.Number },
          { key: 'binanceBid', title: 'Binance Ask', dataType: DataType.Number },
          { key: 'kucoinAsk', title: 'Kucoin Bid', dataType: DataType.Number },
          { key: 'kucoinBid', title: 'Kucoin Ask', dataType: DataType.Number },
          { key: 'binance', title: 'Binance', dataType: DataType.Number },
          { key: 'kucoin', title: 'Kucoin', dataType: DataType.Number },
        ]}
        data={Object.values(table)}
        editingMode={EditingMode.Cell}
        rowKeyField={'id'}
        sortingMode={SortingMode.Single}
      />
    </div>
  );
};