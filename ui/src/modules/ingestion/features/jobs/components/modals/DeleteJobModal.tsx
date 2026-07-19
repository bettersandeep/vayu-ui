import { WarningIcon } from "@phosphor-icons/react"
import { Button, Checkbox, Modal } from "antd"
import { isAxiosError } from "axios"
import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"

import {
	isPostgresCDCSource,
	REPLICATION_SLOT_DROP_WARNING,
} from "@/modules/ingestion/common/utils"
import { useJobStore } from "@/modules/ingestion/features/jobs/stores"

import { useDeleteJob, useJobDetailsFresh } from "../../hooks"

const DeleteJobModal = ({
	fromJobSettings = false,
}: {
	fromJobSettings?: boolean
}) => {
	const { showDeleteJobModal, setShowDeleteJobModal, selectedJobId } =
		useJobStore()
	const { mutateAsync: deleteJob } = useDeleteJob()
	const navigate = useNavigate()
	const [deleteReplicationSlot, setDeleteReplicationSlot] = useState(false)
	const [sharedSlotError, setSharedSlotError] = useState<string | null>(null)

	// The jobs list omits source configs, so fetch the job detail to know
	// whether its source is postgres with CDC (i.e. owns a replication slot).
	const { data: job } = useJobDetailsFresh(
		showDeleteJobModal && selectedJobId ? selectedJobId : undefined,
	)
	const showSlotOption = isPostgresCDCSource(
		job?.source?.type,
		job?.source?.config,
	)

	useEffect(() => {
		if (showDeleteJobModal) {
			setDeleteReplicationSlot(false)
			setSharedSlotError(null)
		}
	}, [showDeleteJobModal])

	const handleDelete = async () => {
		if (!selectedJobId) {
			setShowDeleteJobModal(false)
			return
		}
		try {
			await deleteJob({
				jobId: parseInt(selectedJobId, 10),
				deleteReplicationSlot: showSlotOption && deleteReplicationSlot,
			})
			setShowDeleteJobModal(false)
			if (fromJobSettings) {
				navigate("/jobs")
			}
		} catch (error) {
			if (isAxiosError(error) && error.response?.status === 409) {
				// Shared replication slot: server names the blocking jobs.
				const message = (error.response.data as { message?: string })?.message
				setSharedSlotError(
					message ||
						"The replication slot is used by other jobs. Uncheck the option or delete those jobs first.",
				)
			} else {
				// Error toast is shown by the API interceptor.
				setShowDeleteJobModal(false)
			}
		}
	}

	return (
		<Modal
			open={showDeleteJobModal}
			footer={null}
			closable={false}
			centered
		>
			<div className="flex w-full flex-col items-center justify-center gap-8">
				<WarningIcon
					className="size-16 text-danger"
					weight="fill"
				/>

				<div className="text-center text-xl font-medium text-gray-950">
					Are you sure you want to delete this job?
				</div>

				{showSlotOption && (
					<div className="flex w-full flex-col gap-2">
						<Checkbox
							checked={deleteReplicationSlot}
							onChange={e => {
								setDeleteReplicationSlot(e.target.checked)
								setSharedSlotError(null)
							}}
						>
							Also delete the replication slot used by this job
						</Checkbox>
						{deleteReplicationSlot && (
							<div className="text-xs text-amber-600">
								{REPLICATION_SLOT_DROP_WARNING}
							</div>
						)}
					</div>
				)}

				{sharedSlotError && (
					<div className="w-full rounded-md border border-danger bg-danger-light px-3 py-2 text-sm text-danger">
						{sharedSlotError}
					</div>
				)}

				<div className="flex w-full justify-end gap-4">
					<Button
						type="primary"
						danger
						onClick={handleDelete}
					>
						Delete
					</Button>
					<Button
						type="default"
						onClick={() => setShowDeleteJobModal(false)}
					>
						Cancel
					</Button>
				</div>
			</div>
		</Modal>
	)
}

export default DeleteJobModal
