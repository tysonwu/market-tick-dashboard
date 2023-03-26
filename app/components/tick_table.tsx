import { useState, useEffect } from 'react';
// import ReactApexChart from 'reac/t-apexcharts';
// import { ApexOptions } from 'apexcharts';

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
  'DOT-USDT',
  'TRX-USDT',
  'SHIB-USDT',
  'LTC-USDT',
  'AVAX-USDT',
  'ETH-BTC',
  'BNB-BTC',
  'XRP-BTC',
  'ADA-BTC',
  'DOGE-BTC',
  'MATIC-BTC',
  'SOL-BTC',
  'DOT-BTC',
  'TRX-BTC',
  'LTC-BTC',
  'AVAX-BTC',
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
  const [table, setTable] = useState<PriceTable>(defaultPriceTable);

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

        setTable((prevTable) => {
          var newTable = structuredClone(prevTable);
          if (type === 'ticks') {
            if (exchange == 'binance') { newTable[symbol].binance = msg['Price']; };
            if (exchange == 'kucoin') { newTable[symbol].kucoin = msg['Price']; };
          };
          if (type === 'bidAskTicks') {
            if (exchange === 'binance') { newTable[symbol].binanceBid = msg['Bid']; newTable[symbol].binanceAsk = msg['Ask']; }
            if (exchange === 'kucoin') { newTable[symbol].kucoinBid = msg['Bid']; newTable[symbol].kucoinAsk = msg['Ask']; }
          };
          return newTable;
        });
      };
      // called when cleaning up
      return () => {
        ws.close();
      }
    },
    []
  );

  // const chartOptions: ApexOptions = {
  //   chart: {
  //     type: 'boxPlot',
  //     // height: 350
  //   },
  //   // title: {
  //   //   text: 'Horizontal BoxPlot Chart',
  //   //   align: 'left'
  //   // },
  //   plotOptions: {
  //     bar: {
  //       horizontal: true,
  //       barHeight: '50%'
  //     },
  //     boxPlot: {
  //       colors: {
  //         upper: '#e9ecef',
  //         lower: '#f8f9fa'
  //       }
  //     }
  //   },
  //   stroke: {
  //     colors: ['#6c757d']
  //   }
  // };

  return (
    <div className='cards'>
      {
        symbols.map((symbol) => (
          <div className='card' key={symbol}>
            <h2>{symbol}</h2>
            <div className='exchanges'>
              <div className='exchange' id='binance'>
                <p className='exchange-name'>Binance</p>
                <p className='bid'>{table[symbol].binanceAsk ? table[symbol].binanceAsk : "-"}</p>
                <p className='last-traded'>{table[symbol].binance ? table[symbol].binance : "-"}</p>
                <p className='ask'>{table[symbol].binanceBid ? table[symbol].binanceBid : "-"}</p>
              </div>
              <div className='exchange' id='kucoin'>
                <p className='exchange-name'>Kucoin</p>
                <p className='bid'>{table[symbol].kucoinAsk ? table[symbol].kucoinAsk : "-"}</p>
                <p className='last-traded'>{table[symbol].kucoin ? table[symbol].kucoin : "-"}</p>
                <p className='ask'>{table[symbol].kucoinBid ? table[symbol].kucoinBid : "-"}</p>
              </div>

            </div>
            {/* react-apexchart */}
            {/* <ReactApexChart
              series={[{
                data: [
                  {x: 'Binance', y: [table[symbol].binanceBid, table[symbol].binanceAsk]},
                  {x: 'Kucoin',  y: [table[symbol].kucoinBid, table[symbol].kucoinAsk]},
                ]
              }]}
              options={chartOptions}
              type='rangeBar'
            /> */}
          </div>
        ))
      }
    </div>
  );
};
