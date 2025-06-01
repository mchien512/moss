import {createConnectTransport} from '@connectrpc/connect-web';
import {createClient} from '@connectrpc/connect';
import {
    CreateEntryRequestSchema,
    EntryService,
    GrowthStage
} from '../genproto/protobuf/entry/entry_pb';
import {create} from "@bufbuild/protobuf";

const transport = createConnectTransport({
  baseUrl: process.env.NODE_ENV === 'production'
    ? (process.env.REACT_APP_API_URL || 'https://your-api.com')
    : 'http://localhost:8080',
});

export const entryClient = createClient(EntryService, transport);

// Simple health check using create entry
// export const healthCheck = async (): Promise<boolean> => {
//   try {
//       const request = create(CreateEntryRequestSchema, {
//           title: "Health Check Entry",
//           content: "Testing connection to backend",
//           growthStage: GrowthStage.SEED,
//       });
//
//
//       const response = await entryClient.createEntry(request);
//     console.log('✅ Health check successful - created entry:', response.entry?.id);
//     return true;
//   } catch (error) {
//     console.error('❌ Health check failed:', error);
//     return false;
//   }
// };