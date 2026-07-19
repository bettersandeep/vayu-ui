import { useMutation } from "@tanstack/react-query"

import { notificationService } from "@/core/notifications"
import type {
	EntityBase,
	EntityTestRequest,
} from "@/modules/ingestion/common/types"

import { sourceKeys } from "../../constants/queryKeys"
import { sourceService } from "../../services"

export const useCreateSource = () => {
	return useMutation({
		mutationKey: sourceKeys.all(),
		mutationFn: (source: EntityBase) => sourceService.createSource(source),
	})
}

export const useUpdateSource = (id: string) => {
	return useMutation({
		mutationKey: sourceKeys.all(),
		mutationFn: (source: EntityBase) => sourceService.updateSource(id, source),
	})
}

export const useDeleteSource = () => {
	return useMutation({
		mutationKey: sourceKeys.all(),
		mutationFn: ({
			id,
			deleteReplicationSlot,
		}: {
			id: string
			deleteReplicationSlot?: boolean
		}) => sourceService.deleteSource(id, deleteReplicationSlot),
		onSuccess: data => {
			// Source deletion succeeded but the replication slot drop did not.
			if (data?.replication_slot_warning) {
				notificationService.warning(data.replication_slot_warning)
			}
		},
	})
}

export const useTestSourceConnection = () => {
	return useMutation({
		mutationFn: ({
			source,
			existing = false,
		}: {
			source: EntityTestRequest
			existing?: boolean
		}) => sourceService.testSourceConnection(source, existing),
	})
}
