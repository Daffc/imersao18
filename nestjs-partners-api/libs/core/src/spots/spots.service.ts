import { HttpCode, Injectable } from '@nestjs/common';
import { CreateSpotDto } from './dto/create-spot.dto';
import { UpdateSpotDto } from './dto/update-spot.dto';
import { PrismaService } from '../prisma/prisma.service';
import { SpotStatus } from '@prisma/client';

@Injectable()
export class SpotsService {
  constructor(private prismaService: PrismaService) { }
  async create(createSpotDto: CreateSpotDto & { eventId: string }) {

    const event = await this.prismaService.event.findFirst({
      where: { id: createSpotDto.eventId }
    });

    if (!event) {
      throw new Error('Event not found.')
    }

    return await this.prismaService.spot.create({
      data: {
        ...createSpotDto,
        status: SpotStatus.available,
      }
    })
  }

  async findAll(eventId: string) {
    return await this.prismaService.spot.findMany({
      where: { eventId: eventId }
    })
  }

  async findOne(eventId: string, spotId: string) {
    return await this.prismaService.spot.findFirst({
      where: {
        id: spotId,
        eventId: eventId
      }
    })
  }

  async update(eventId: string, spotId: string, updateSpotDto: UpdateSpotDto) {
    return await this.prismaService.spot.update({
      data: updateSpotDto,
      where: {
        id: spotId,
        eventId: eventId
      }
    })
  }

  async remove(eventId: string, spotId: string) {
    return await this.prismaService.spot.delete({
      where: {
        id: spotId,
        eventId: eventId
      }
    })
  }
}
