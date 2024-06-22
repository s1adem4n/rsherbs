import Pocketbase, { type RecordService } from 'pocketbase';

interface Base {
	id: string;
	created: string;
	updated: string;
}

export interface Plant extends Base {
	name: string;
	latin: string;
	description: string;
	image: string;
}

interface TypedPocketbase extends Pocketbase {
	collection(idOrName: string): RecordService;
	collection(idOrName: 'plants'): RecordService<Plant>;
}

export const BASE_URL = 'http://localhost:8090' as const;

const pb = new Pocketbase(BASE_URL) as TypedPocketbase;

export default pb;
