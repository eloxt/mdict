import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Dictionary } from '~/main'

interface DictStore {
    dictionary: Dictionary | null
    word: string | null
    setDictionary: (dictionary: Dictionary) => void
    setWord: (word: string) => void
}

export const useDictStore = create<DictStore>()(
    persist(
        (set) => ({
            dictionary: null,
            word: null,
            setDictionary: (dictionary: Dictionary) => set({ dictionary }),
            setWord: (word: string) => set({ word }),
        }),
        {
            name: 'dict-store',
            partialize: (state) => ({ dictionary: state.dictionary }),
        }
    )
)
